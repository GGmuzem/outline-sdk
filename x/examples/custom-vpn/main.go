package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"fyne.io/fyne/v2/app"
	"golang.getoutline.org/sdk/x/configurl"
	"golang.getoutline.org/sdk/x/httpproxy"
)

var (
	proxyServer      *http.Server
	currentProxyAddr string
)

func startVPN(config string) error {
	dialer, err := configurl.NewDefaultProviders().NewStreamDialer(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create dialer: %w", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	currentProxyAddr = listener.Addr().String()
	host, port, _ := net.SplitHostPort(currentProxyAddr)

	proxyServer = &http.Server{
		Handler: httpproxy.NewProxyHandler(dialer),
	}

	go func() {
		if err := proxyServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("Proxy server error: %v\n", err)
		}
	}()

	if err := setSystemProxy(host, port); err != nil {
		proxyServer.Close()
		return fmt.Errorf("failed to set system proxy: %w", err)
	}

	return nil
}

func stopVPN() error {
	if proxyServer != nil {
		proxyServer.Close()
		proxyServer = nil
	}
	return unsetSystemProxy()
}

func main() {
	transportConfig := flag.String("transport", "", "Transport config (ss://...)")
	flag.Parse()

	myApp := app.New()
	win := setupGUI(myApp)

	log.Printf("Starting Dr. Frake VPN with config: %s\n", *transportConfig)

	win.ShowAndRun()

	// Ensure proxy is unset on exit
	stopVPN()
}
