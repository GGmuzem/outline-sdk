package main

import (
	"drfrake-backend/outline"
)

// OutlineProvider implements VPNProvider using the Outline Server API.
type OutlineProvider struct {
	client *outline.Client
}

// NewOutlineProvider creates a provider backed by an Outline Server.
func NewOutlineProvider(apiURL, certSHA256 string) *OutlineProvider {
	return &OutlineProvider{
		client: outline.NewClient(apiURL, certSHA256),
	}
}

func (p *OutlineProvider) CreateKey(userID string) (string, string, error) {
	key, err := p.client.CreateKey()
	if err != nil {
		return "", "", err
	}
	// Set name for tracking
	p.client.SetName(key.ID, "user-"+userID)
	return key.ID, key.AccessURL, nil
}

func (p *OutlineProvider) DeleteKey(keyID string) error {
	return p.client.DeleteKey(keyID)
}

func (p *OutlineProvider) GetKeys() ([]VPNKey, error) {
	keys, err := p.client.GetKeys()
	if err != nil {
		return nil, err
	}
	var result []VPNKey
	for _, k := range keys {
		result = append(result, VPNKey{
			ID:        k.ID,
			Name:      k.Name,
			AccessURL: k.AccessURL,
		})
	}
	return result, nil
}

func (p *OutlineProvider) SetName(keyID string, name string) error {
	return p.client.SetName(keyID, name)
}
