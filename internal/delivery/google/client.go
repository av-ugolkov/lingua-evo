package google

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	TokenInfo   = "https://oauth2.googleapis.com/tokeninfo?id_token=%s"
	RemoveTokem = "https://oauth2.googleapis.com/revoke"
)

var googleConfig *oauth2.Config

func InitClient(cfg *config.Google) {
	secret, err := os.ReadFile(cfg.SecretPath)
	if err != nil {
		panic(err)
	}

	googleConfig, err = google.ConfigFromJSON(
		secret,
		"openid",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile")
	if err != nil {
		panic(err)
	}
}

func GetAuthUrl() string {
	return googleConfig.AuthCodeURL(uuid.New().String(), oauth2.AccessTypeOffline)
}

func GetTokenByCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := googleConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("google.GetTokenByCode: %w", err)
	}

	return token, nil
}

func GetProfile(ctx context.Context, token *oauth2.Token) (*GoogleProfile, error) {
	httpClient := googleConfig.Client(ctx, token)
	resp, err := httpClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("google.GetProfile: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("google.GetProfile: %w", err)
	}

	profile := new(GoogleProfile)
	err = jsoniter.Unmarshal(data, profile)
	if err != nil {
		return nil, fmt.Errorf("google.GetProfile: %w", err)
	}

	return profile, nil
}

func RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	tokenURL := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("client_id", googleConfig.ClientID)
	data.Set("client_secret", googleConfig.ClientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to refresh token, status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := jsoniter.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format: %v", result)
	}

	return accessToken, nil
}

func VerifyGoogleToken(idToken string, clientID string) (*GoogleTokenInfo, error) {
	resp, err := http.Get(fmt.Sprintf(TokenInfo, idToken))
	if err != nil {
		return nil, fmt.Errorf("delivery.google.VerifyGoogleToken: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("delivery.google.VerifyGoogleToken: the checking token was failed")
	}

	var tokenInfo GoogleTokenInfo
	if err := jsoniter.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("delivery.google.VerifyGoogleToken: %v", err)
	}

	if tokenInfo.Audience != clientID {
		return nil, fmt.Errorf("delivery.google.VerifyGoogleToken: mismatch client_id: %s", tokenInfo.Audience)
	}

	if time.Now().Unix() > int64(tokenInfo.ExpiresIn) {
		return nil, fmt.Errorf("delivery.google.VerifyGoogleToken: token is expired")
	}

	return &tokenInfo, nil
}

func RevokeGoogleToken(token string) error {
	data := url.Values{}
	data.Set("client_id", googleConfig.ClientID)
	data.Set("client_secret", googleConfig.ClientSecret)
	data.Set("token", token)

	req, err := http.NewRequest("POST", RemoveTokem, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("delivery.google.RevokeGoogleToken: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data.Encode())))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("delivery.google.RevokeGoogleToken: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delivery.google.RevokeGoogleToken: can't revoke token")
	}

	return nil
}
