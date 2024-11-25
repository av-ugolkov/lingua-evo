package google

type (
	GoogleProfile struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
		VerifiedEmail bool   `json:"verified_email"`
	}

	GoogleTokenInfo struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		ExpiresIn     int    `json:"exp"`
		Audience      string `json:"aud"` // Ваш client_id
		Issuer        string `json:"iss"`
	}
)
