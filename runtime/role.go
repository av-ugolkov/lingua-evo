package runtime

type Role string

const (
	User  Role = "user"
	Guest Role = "guest"
	Admin Role = "admin"
)

func (r Role) IsAdmin() bool {
	return r == Admin
}
