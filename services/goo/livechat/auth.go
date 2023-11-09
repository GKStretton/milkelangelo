package livechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gkstretton/dark/services/goo/keyvalue"
)

const (
	tokenEndpoint   = "https://id.twitch.tv/oauth2/token"
	refreshTokenKey = "TWITCH_REFRESH_TOKEN"
)

// TokenResponse represents the response from the Twitch token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // The remaining time before the token expires, in seconds
}

// RefreshAccessToken uses the refresh token to get a new access token
func RefreshAccessToken(clientID, clientSecret, refreshToken string) (*TokenResponse, error) {
	body := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&client_id=%s&client_secret=%s", refreshToken, clientID, clientSecret)
	req, err := http.NewRequest("POST", tokenEndpoint, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: time.Second * time.Duration(10),
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error %s refreshing token: %v", resp.Status, string(respBody))
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(respBody, &tokenResponse)
	if err != nil {
		return nil, err
	}

	fmt.Printf("successfully refreshed twitch access token\n")

	return &tokenResponse, nil
}

func getRefreshToken() string {
	return string(keyvalue.Get(refreshTokenKey))
}

func setRefreshToken(token string) {
	keyvalue.Set(refreshTokenKey, []byte(token))
}
