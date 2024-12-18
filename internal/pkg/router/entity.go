package router

import "github.com/av-ugolkov/lingua-evo/runtime"

const (
	Language       = "language"
	RefreshToken   = "refresh_token"
	CookiePathAuth = "/auth"
)
const (
	HeaderAuthorization = "Authorization"
	HeaderFingerprint   = "Fingerprint"
)

const (
	AuthTypeNone   = runtime.EmptyString
	AuthTypeBearer = "Bearer"
	AuthTypeBasic  = "Basic"
)
