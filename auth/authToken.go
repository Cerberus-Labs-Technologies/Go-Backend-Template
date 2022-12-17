package auth

import (
	"necross.it/backend/database"
	"necross.it/backend/util"
)

type Token struct {
	Id        int              `json:"id" db:"id"`
	UserId    int              `json:"user_id" db:"user_id"`
	Scope     int              `json:"scope" db:"scope"`
	Token     string           `json:"token" db:"token"`
	CreatedAt util.TimeStamp   `json:"createdAt" db:"createdAt"`
	UpdatedAt util.TimeStamp   `json:"updatedAt" db:"updatedAt"`
	ExpiresAt util.TimeStamp   `json:"expiresAt" db:"expiresAt"`
	Active    database.IntBool `json:"active" db:"active"`
}
