package model

type RoleID string

const (
	RoleAdmin RoleID = "admin"
	RoleUser  RoleID = "user"
)

type Role struct {
	ID   RoleID
	Name string
}

func (r *Role) IsAdmin() bool {
	return r.ID == RoleAdmin
}
