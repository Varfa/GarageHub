package models

import "time"

type Role struct {
	ID int64

	Code        string
	Name        string
	Description *string

	IsSystem bool
	IsActive bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
type RoleListItem struct {
	ID               int64
	Code             string
	Name             string
	Description      string
	IsSystem         bool
	IsActive         bool
	PermissionsCount int
}
type RoleDetails struct {
	ID          int64
	Code        string
	Name        string
	Description string
	IsSystem    bool
	IsActive    bool
}
type RolePermissionItem struct {
	ID          int64
	Code        string
	Module      string
	Name        string
	Description string
	Assigned    bool
}
type PermissionGroup struct {
	Module      string
	Permissions []RolePermissionItem
}
