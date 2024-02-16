package exchange

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

type ContentType string

const (
	ContentTypeJSON ContentType = "application/json"
)
