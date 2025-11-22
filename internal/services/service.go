package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func getOAuthToken(consumerKey, consumerSecret string) (string, error) {
	tokenURL := "https://api.login.yahoo.com/oauth/v2/get_request_token"

	auth := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))

	body := []byte("grant_type=client_credentials")

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token: %s", string(respBytes))
	}

	var tokenResp OAuthTokenResponse
	if err := json.Unmarshal(respBytes, &tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func getData(token string) (string, error) {
	apiURL := "https://api.fantasydata.net/v3/nfl/scores/json/Players" // This URL is not real; replace with actual API endpoint

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get data: %s", string(respBytes))
	}

	return string(respBytes), nil
}
