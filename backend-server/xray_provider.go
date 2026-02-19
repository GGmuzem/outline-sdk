package main

import (
	"encoding/json"
	"fmt"
	"log"

	"drfrake-backend/xray"

	"github.com/google/uuid"
)

// XrayProvider implements VPNProvider using 3X-UI panel API.
type XrayProvider struct {
	client     *xray.Client
	inboundID  int
	serverHost string // Public IP/hostname of the VPN server
	settings   XrayServerSettings
}

// XrayServerSettings holds server-specific VLESS+Reality parameters.
type XrayServerSettings struct {
	Port        int    `json:"port"`
	Flow        string `json:"flow"`
	Security    string `json:"security"`    // "reality"
	SNI         string `json:"sni"`         // e.g. "google.com"
	Fingerprint string `json:"fingerprint"` // e.g. "chrome"
	PublicKey   string `json:"public_key"`
	ShortID     string `json:"short_id"`
	SpiderX     string `json:"spider_x"`
}

// NewXrayProvider creates a provider backed by a 3X-UI panel.
func NewXrayProvider(panelURL, username, password string, inboundID int, serverHost string, settingsJSON string) *XrayProvider {
	var settings XrayServerSettings
	if err := json.Unmarshal([]byte(settingsJSON), &settings); err != nil {
		log.Printf("Warning: failed to parse xray settings: %v", err)
		settings = XrayServerSettings{
			Port:        443,
			Flow:        "xtls-rprx-vision",
			Security:    "reality",
			SNI:         "google.com",
			Fingerprint: "chrome",
		}
	}

	return &XrayProvider{
		client:     xray.NewClient(panelURL, username, password),
		inboundID:  inboundID,
		serverHost: serverHost,
		settings:   settings,
	}
}

func (p *XrayProvider) CreateKey(userID string) (string, string, error) {
	clientUUID := uuid.New().String()
	email := fmt.Sprintf("user-%s", userID)

	if err := p.client.AddClient(p.inboundID, clientUUID, email); err != nil {
		return "", "", fmt.Errorf("failed to create xray client: %w", err)
	}

	// Build VLESS URI
	vlessURI := xray.BuildVLESSURI(xray.VLESSConfig{
		UUID:        clientUUID,
		Host:        p.serverHost,
		Port:        p.settings.Port,
		Flow:        p.settings.Flow,
		Security:    p.settings.Security,
		SNI:         p.settings.SNI,
		Fingerprint: p.settings.Fingerprint,
		PublicKey:   p.settings.PublicKey,
		ShortID:     p.settings.ShortID,
		SpiderX:     p.settings.SpiderX,
	})

	return clientUUID, vlessURI, nil
}

func (p *XrayProvider) DeleteKey(keyID string) error {
	return p.client.RemoveClient(p.inboundID, keyID)
}

func (p *XrayProvider) GetKeys() ([]VPNKey, error) {
	clients, err := p.client.GetClients(p.inboundID)
	if err != nil {
		return nil, err
	}

	var keys []VPNKey
	for _, c := range clients {
		vlessURI := xray.BuildVLESSURI(xray.VLESSConfig{
			UUID:        c.ID,
			Host:        p.serverHost,
			Port:        p.settings.Port,
			Flow:        p.settings.Flow,
			Security:    p.settings.Security,
			SNI:         p.settings.SNI,
			Fingerprint: p.settings.Fingerprint,
			PublicKey:   p.settings.PublicKey,
			ShortID:     p.settings.ShortID,
			SpiderX:     p.settings.SpiderX,
		})

		keys = append(keys, VPNKey{
			ID:        c.ID,
			Name:      c.Email,
			AccessURL: vlessURI,
		})
	}
	return keys, nil
}

func (p *XrayProvider) SetName(keyID string, name string) error {
	// 3X-UI uses email as identifier; name change not easily supported via API
	// This is a no-op for now
	return nil
}
