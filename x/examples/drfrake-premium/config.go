package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	BackendURL string `json:"backend_url"`
}

type ServerConfig struct {
	ID        string `json:"id"`
	Country   string `json:"country"`
	Flag      string `json:"flag"`
	City      string `json:"city"`
	Config    string `json:"config"` // SS URI
	IsPremium bool   `json:"isPremium"`
	IsDefault bool   `json:"isDefault"`
}

// Server is the struct exposed to the frontend
type Server struct {
	ID        string `json:"id"`
	Country   string `json:"country"`
	City      string `json:"city"`
	Flag      string `json:"flag"`
	Config    string `json:"config"`
	IsPremium bool   `json:"isPremium"`
	Latency   int    `json:"latency"`
}

func GetConfigDir() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, "DrFrakeVPN")
}

func LoadConfig() (*Config, error) {
	path := filepath.Join(GetConfigDir(), "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		// Create default config if missing
		cfg := &Config{
			BackendURL: "http://localhost:8080",
		}
		SaveConfig(cfg)
		return cfg, nil
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	path := filepath.Join(GetConfigDir(), "config.json")
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(path, data, 0644)
}

func LoadServers() ([]ServerConfig, error) {
	path := filepath.Join(GetConfigDir(), "servers.json")
	data, err := os.ReadFile(path)
	if err != nil {
		// Return default fallback servers + write them to file
		defaults := GetDefaultServers()
		SaveServers(defaults)
		return defaults, nil
	}

	var servers []ServerConfig
	if err := json.Unmarshal(data, &servers); err != nil {
		return GetDefaultServers(), nil
	}
	return servers, nil
}

func SaveServers(servers []ServerConfig) error {
	path := filepath.Join(GetConfigDir(), "servers.json")
	data, _ := json.MarshalIndent(servers, "", "  ")
	return os.WriteFile(path, data, 0644)
}

func GetDefaultServers() []ServerConfig {
	// These match the old hardcoded ones
	return []ServerConfig{
		{ID: "us-1", Country: "USA", Flag: "🇺🇸", City: "New York", Config: "ss://Y2hhY2hhMjAtaWV0Zi1wb2x5MTMwNTpwYXNzd29yZA@127.0.0.1:8388#USA-Server", IsDefault: true},
		{ID: "nl-1", Country: "Netherlands", Flag: "🇳🇱", City: "Amsterdam", Config: "ss://Y2hhY2hhMjAtaWV0Zi1wb2x5MTMwNTpwYXNzd29yZA@127.0.0.1:8389#NL-Server"},
		{ID: "jp-1", Country: "Japan", Flag: "🇯🇵", City: "Tokyo", Config: "ss://Y2hhY2hhMjAtaWV0Zi1wb2x5MTMwNTpwYXNzd29yZA@127.0.0.1:8390#JP-Premium", IsPremium: true},
		{ID: "de-1", Country: "Germany", Flag: "🇩🇪", City: "Frankfurt", Config: "ss://Y2hhY2hhMjAtaWV0Zi1wb2x5MTMwNTpwYXNzd29yZA@127.0.0.1:8391#DE-Premium", IsPremium: true},
	}
}
