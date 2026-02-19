package outline

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	APIURL     string
	CertSHA256 string
	httpClient *http.Client
}

type AccessKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
	Method    string `json:"method"`
	AccessURL string `json:"accessUrl"`
}

func NewClient(apiURL, certSHA256 string) *Client {
	// Outline uses self-signed certs, so we need to skip verification or pin cert
	// For simplicity in this implementation, we'll skip verify but real prod should pin cert
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		APIURL:     apiURL,
		CertSHA256: certSHA256,
		httpClient: &http.Client{Transport: tr, Timeout: 10 * time.Second},
	}
}

func (c *Client) CreateKey() (*AccessKey, error) {
	resp, err := c.httpClient.Post(c.APIURL+"/access-keys", "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("outline api error: %d", resp.StatusCode)
	}

	var key AccessKey
	if err := json.NewDecoder(resp.Body).Decode(&key); err != nil {
		return nil, err
	}
	return &key, nil
}

func (c *Client) GetKeys() ([]AccessKey, error) {
	resp, err := c.httpClient.Get(c.APIURL + "/access-keys")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("outline api error: %d", resp.StatusCode)
	}

	var result struct {
		AccessKeys []AccessKey `json:"accessKeys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.AccessKeys, nil
}

func (c *Client) DeleteKey(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/access-keys/%s", c.APIURL, id), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("outline api error: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) SetName(id, name string) error {
	payload := map[string]string{"name": name}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/access-keys/%s/name", c.APIURL, id), strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("outline api error: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) SetDataLimit(id string, bytes int64) error {
	url := fmt.Sprintf("%s/access-keys/%s/data-limit", c.APIURL, id)

	var payload interface{}
	if bytes > 0 {
		payload = map[string]interface{}{
			"limit": map[string]int64{"bytes": bytes},
		}
	} else {
		// To remove limit, we send DELETE request to data-limit endpoint
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return err
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return nil
	}

	data, _ := json.Marshal(payload)
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("outline api error: %d", resp.StatusCode)
	}
	return nil
}
