package sappress

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type authTransport struct {
	token string
}

type LoginResponse struct {
	DeviceStatus string `json:"device_status"`
	Token        string `json:"token"`
	UserKey      string `json:"user_key"`
}

// CheckToken validates token and refreshes if needed
func CheckToken() {
	config, _ := LoadConfig()
	client := NewAuthenticatedClient(config.Token)
	resp, err := client.Get(AccountListURL)
	if err != nil {
		log.Fatal("ailed to call AccountListURL:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Token expired or invalid, Re-authenticating...")

		newToken, userKey, err := getNewToken()
		if err != nil {
			log.Println("Failed to refresh token:", err)
			return
		}
		if newToken != "" {
			UpdateConfigField("token", newToken)
			UpdateConfigField("userkey", userKey)
			log.Println("logged in and updated token successfully")
		}
		return
	}
	log.Println("Token is valid")
}

// returns a fresh token and the userkey
func getNewToken() (string, string, error) {
	conf, err := LoadConfig()
	if err != nil {
		return "", "", err
	}

	payload := map[string]interface{}{
		"email":                     conf.Email,
		"password":                  conf.Password,
		"app_version":               "0",
		"device_id":                 "b3f1a8c7-9e2d-4a52-891b-36f10ae487d9",
		"os_type":                   "android",
		"os_version":                "11",
		"device_name_manufacturer": "Pixel 4a",
		"device_name_user":         "Pixel 4a",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", LoginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "okhttp/4.12.0")
	req.Header.Set("x-project", "sap-press")
	req.Header.Set("Host", "eba.sap-press.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("login failed: %s", string(body))
	}

	var loginResp LoginResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", "", err
	}

	return loginResp.Token, loginResp.UserKey, nil
}

// NewAuthenticatedClient creates an http.Client with auth headers
func NewAuthenticatedClient(token string) *http.Client {
	return &http.Client{
		Transport: &authTransport{token: token},
	}
}

// Injects auth headers into every request
func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Token "+a.token)
	req.Header.Set("x-project", "sap-press")
	req.Header.Set("Host", "eba.sap-press.com")
	req.Header.Set("User-Agent", "okhttp/4.12.0")
	req.Header.Set("Accept-Encoding", "gzip")
	return http.DefaultTransport.RoundTrip(req)
}
