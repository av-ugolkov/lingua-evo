package ginext

import "github.com/av-ugolkov/lingua-evo/runtime"

const (
	Language     = "language"
	RefreshToken = "refresh_token"
)
const (
	Authorization = "Authorization"
	Fingerprint   = "Fingerprint"
)

const (
	AuthTypeNone   = runtime.EmptyString
	AuthTypeBearer = "Bearer"
	AuthTypeBasic  = "Basic"
)
