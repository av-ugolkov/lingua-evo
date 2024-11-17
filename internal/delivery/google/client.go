package google

import (
	"context"
	_ "embed"
	"fmt"
	"io"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	jsoniter "github.com/json-iterator/go"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

//go:embed google_client_secret.json
var googleClientSecret string

var googleConfig *oauth2.Config

func init() {
	var err error
	googleConfig, err = google.ConfigFromJSON(
		[]byte(googleClientSecret),
		"openid",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile")
	if err != nil {
		panic(err)
	}
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

	var profile *GoogleProfile
	err = jsoniter.Unmarshal(data, profile)
	if err != nil {
		return nil, fmt.Errorf("google.GetProfile: %w", err)
	}

	return profile, nil
}

func ParseGoogleOauthToken(ctx context.Context, token string) (any, error) {
	payload, err := idtoken.Validate(ctx, token, config.GetConfig().Google.ClientID)
	if err != nil {
		return nil, fmt.Errorf("token.ParseGoogleOauthToken: %w", err)
	}

	return payload, nil
}
