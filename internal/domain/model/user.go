package model

import "time"

type User struct {
	ID             int
	RoleID         RoleID
	Username       string
	Email          string
	Password       string
	EmailVerified  bool
	SessionVersion int64
	CreatedAt      time.Time
}
