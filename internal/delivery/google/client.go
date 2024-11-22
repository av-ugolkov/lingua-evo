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
	// Формируем URL запроса
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)

	// Отправляем запрос к Google
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("не удалось отправить запрос: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не удалось проверить токен, статус: %v", resp.StatusCode)
	}

	// Парсим ответ
	var tokenInfo GoogleTokenInfo
	if err := jsoniter.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	// Проверяем клиента и срок действия токена
	if tokenInfo.Audience != clientID {
		return nil, fmt.Errorf("несовпадение client_id: %s", tokenInfo.Audience)
	}

	if time.Now().Unix() > int64(tokenInfo.ExpiresIn) {
		return nil, fmt.Errorf("токен истёк")
	}

	return &tokenInfo, nil
}

func RevokeGoogleToken(accessToken string) error {
	tokenURL := "https://oauth2.googleapis.com/revoke"

	data := url.Values{}
	data.Set("client_id", googleConfig.ClientID)
	data.Set("client_secret", googleConfig.ClientSecret)
	data.Set("token", accessToken)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data.Encode())))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("не удалось удалить токен, статус: %v", resp.StatusCode)
	}

	return nil
}

// func ParseGoogleOauthToken(ctx context.Context, token string) (any, error) {
// 	payload, err := idtoken.Validate(ctx, token, config.GetConfig().Google.ClientID)
// 	if err != nil {
// 		return nil, fmt.Errorf("token.ParseGoogleOauthToken: %w", err)
// 	}

// 	return payload, nil
// }
