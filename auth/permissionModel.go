package auth

import "necross.it/backend/database"

type Group struct {
	ID        int                `json:"id" db:"id"`
	Name      string             `json:"name" db:"name"`
	CreatedAt database.TimeStamp `json:"createdAt" db:"created_at"`
	UpdatedAt database.TimeStamp `json:"updatedAt" db:"updated_at"`
}

type GroupPermission struct {
	ID           int                `json:"id" db:"id"`
	GroupId      int                `json:"groupId" db:"group_id"`
	PermissionId int                `json:"permissionId" db:"permission_id"`
	CreatedAt    database.TimeStamp `json:"createdAt" db:"created_at"`
	UpdatedAt    database.TimeStamp `json:"updatedAt" db:"updated_at"`
}

type UserPermission struct {
	ID           int                `json:"id" db:"id"`
	UserId       int                `json:"userId" db:"user_id"`
	PermissionId int                `json:"permissionId" db:"permission_id"`
	CreatedAt    database.TimeStamp `json:"createdAt" db:"created_at"`
	UpdatedAt    database.TimeStamp `json:"updatedAt" db:"updated_at"`
}

type Permission struct {
	ID        int                `json:"id" db:"id"`
	Name      string             `json:"name" db:"name"`
	CreatedAt database.TimeStamp `json:"createdAt" db:"created_at"`
	UpdatedAt database.TimeStamp `json:"updatedAt" db:"updated_at"`
}
