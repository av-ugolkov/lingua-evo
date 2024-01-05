package common

const (
	Language     = "language"
	RefreshToken = "refresh_token"
)

type TypeAuth string

const (
	AuthTypeNone   TypeAuth = ""
	AuthTypeBearer TypeAuth = "Bearer"
	AuthTypeBasic  TypeAuth = "Basic"
)
