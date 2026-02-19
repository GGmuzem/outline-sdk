package core

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.getoutline.org/sdk/x/configurl"
	"golang.getoutline.org/sdk/x/httpproxy"
)

// VPNClient manages the connection
type VPNClient struct {
	proxyServer  *http.Server
	isConnected  bool
	activeConfig string
}

func NewVPNClient() *VPNClient {
	return &VPNClient{}
}

// Connect starts the local proxy and returns the bound address (host:port).
// On mobile, the UI layer (Flutter/Kotlin/Swift) must route traffic to this address.
func (c *VPNClient) Connect(config string) (string, error) {
	if c.isConnected {
		return "", fmt.Errorf("already connected")
	}

	dialer, err := configurl.NewDefaultProviders().NewStreamDialer(context.Background(), config)
	if err != nil {
		return "", fmt.Errorf("failed to create dialer: %w", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("failed to listen: %w", err)
	}

	proxyAddr := listener.Addr().String()

	c.proxyServer = &http.Server{
		Handler: httpproxy.NewProxyHandler(dialer),
	}

	go func() {
		if err := c.proxyServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("Proxy server error: %v\n", err)
		}
	}()

	c.isConnected = true
	c.activeConfig = config

	// Return the address so mobile native layer can use it (VpnService/tun2socks)
	return proxyAddr, nil
}

func (c *VPNClient) Disconnect() error {
	if c.proxyServer != nil {
		c.proxyServer.Close()
		c.proxyServer = nil
	}
	c.isConnected = false
	return nil
}

func (c *VPNClient) IsConnected() bool {
	return c.isConnected
}
