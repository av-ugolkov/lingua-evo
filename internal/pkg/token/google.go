package token

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/config"

	"google.golang.org/api/idtoken"
)

func ParseGoogleOauthToken(ctx context.Context, token string) (any, error) {
	payload, err := idtoken.Validate(ctx, token, config.GetConfig().Google.ClientID)
	if err != nil {
		return nil, fmt.Errorf("token.ParseGoogleOauthToken: %w", err)
	}

	return payload, nil
}
