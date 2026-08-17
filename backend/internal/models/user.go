package models

import "time"

type User struct {
	ID int64

	EmployeeID *int64
	RoleID     *int64

	Email        string
	PasswordHash string

	IsOwner  bool
	IsActive bool

	LastLoginAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
type UserListItem struct {
	ID           int64
	EmployeeID   int64
	FirstName    string
	LastName     string
	PositionName string
	Email        string
	RoleID       *int64
	RoleIDValue  int64
	RoleCode     string
	RoleName     string
	IsActive     bool
	IsOwner      bool
}
