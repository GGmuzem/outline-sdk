package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthClient struct {
	BaseURL string
	Token   string
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{BaseURL: baseURL}
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Plan  string `json:"plan"`
	} `json:"user"`
}

func (c *AuthClient) Login(email, password string) error {
	payload := map[string]string{"email": email, "password": password}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(c.BaseURL+"/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	c.Token = authResp.Token
	return nil
}

func (c *AuthClient) GetServers() ([]string, error) {
	req, _ := http.NewRequest("GET", c.BaseURL+"/servers", nil)
	req.Header.Set("Authorization", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch servers: %s", resp.Status)
	}

	// Assuming servers returns a list of objects with "config" field
	var serverList []struct {
		Config string `json:"config"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&serverList); err != nil {
		return nil, err
	}

	var configs []string
	for _, s := range serverList {
		configs = append(configs, s.Config)
	}
	return configs, nil
}
