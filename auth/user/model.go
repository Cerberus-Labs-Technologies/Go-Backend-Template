package user

import (
	"math/rand"
	"necross.it/backend/database"
	"necross.it/backend/util"
)

type User struct {
	Id              int                 `json:"id" db:"id"`
	Name            string              `json:"name" db:"name"`
	Email           string              `json:"email" db:"email"`
	EmailVerifiedAt util.TimeStamp      `json:"email_verified_at" db:"email_verified_at"`
	Password        string              `json:"password" db:"password"`
	Scope           string              `json:"scope" db:"scope"`
	RememberToken   database.JsonString `json:"remember_token" db:"remember_token"`
	CreatedAt       util.TimeStamp      `json:"created_at" db:"created_at"`
	UpdatedAt       util.TimeStamp      `json:"updated_at" db:"updated_at"`
}

func (u *User) CreateAuthToken() string {
	token := GenerateRandomToken()
	return token
}

func (u *User) ConvertToAuthJSON() AuthJSON {
	return AuthJSON{
		Id:              u.Id,
		Name:            u.Name,
		Email:           u.Email,
		EmailVerifiedAt: u.EmailVerifiedAt,
		Scope:           u.Scope,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

type AuthJSON struct {
	Id              int            `json:"id"`
	Name            string         `json:"name"`
	Email           string         `json:"email"`
	EmailVerifiedAt util.TimeStamp `json:"email_verified_at"`
	Scope           string         `json:"scope"`
	CreatedAt       util.TimeStamp `json:"created_at"`
	UpdatedAt       util.TimeStamp `json:"updated_at"`
}

type ForgotPassword struct {
	ID         int            `json:"id" db:"ID"`
	UserID     int            `json:"userId" db:"userId"`
	ResetToken string         `json:"resetToken" db:"resetToken"`
	CreatedAt  util.TimeStamp `json:"createdAt" db:"createdAt"`
}

func GenerateRandomToken() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 991)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
