package auth

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"necross.it/backend/auth/user"
	"necross.it/backend/database"
	"necross.it/backend/util"
	"necross.it/backend/verify"
)

type Service struct {
	Server  database.Server
	User    user.Service
	Verify  verify.Service
	General util.Service
}

func (s *Service) RegisterRoutes(authGroup fiber.Router) {
	authGroup.Post("/login", s.LoginAuth)
	authGroup.Get("/user", s.UserMe)
	authGroup.Post("/logout", s.Logout)
	authGroup.Post("/register", s.RegisterController)
}

func (s *Service) RegisterUserRoutes(userGroup fiber.Router) {
	userGroup.Put("/change-email", s.IsLoggedIn, s.ChangeEmail)
	userGroup.Delete("/delete", s.IsLoggedIn, s.DeleteAccount)
	userGroup.Put("/change-password", s.IsLoggedIn, s.ChangePassword)
	userGroup.Delete("/delete-account", s.IsLoggedIn, s.DeleteAccountConfirm)
	userGroup.Put("/forgot-password", s.ForgotPassword)
	userGroup.Put("/reset-password", s.ResetPassword)
}

// Login password should be encrypted before entry
func (s *Service) Login(email string, password []byte) (user.User, string, error) {
	emptyUser := user.User{}
	if (email == "") || (password == nil) {
		return emptyUser, "", errors.New("Email oder Passwort ist leer")
	}
	if s.UserExists(email) {
		user, err := s.GetUserByEmail(email)
		if bcrypt.CompareHashAndPassword([]byte(user.Password), password) == nil {
			var token = user.CreateAuthToken()
			return user, token, err
		} else {
			return emptyUser, "", errors.New("Passwort oder email sind falsch")
		}
	} else {
		return emptyUser, "", errors.New("User wurde nicht gefunden")
	}
}

func (s *Service) UserExists(email string) bool {
	var user user.User
	rs := !s.General.EntryExists(user, "users", "email", email)
	return rs
}

// UserExistsByEmail check if user exists by email in database
func (s *Service) UserExistsByEmail(email string) bool {
	var user user.User
	err := s.Server.DB.Get(&user, "SELECT * FROM users WHERE email = ?", email)
	return err == nil
}

func (s *Service) RegisterService(registrationForm RegisterForm) (user.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationForm.Password), bcrypt.DefaultCost)
	user := user.User{
		Name:     registrationForm.Name,
		Email:    registrationForm.Email,
		Password: string(hashedPassword),
	}
	err = s.User.Register(user)
	createdUser, err := s.GetUserByEmail(user.Email)
	err = s.Verify.SendVerifyEmail(createdUser)
	return createdUser, err
}

func (s *Service) GetUserByEmail(email string) (user.User, error) {
	var user user.User
	err := s.Server.DB.Get(&user, "SELECT * FROM users WHERE email = ?", email)
	return user, err
}

type FormAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CheckScope func check scope of user
func (s *Service) CheckScope(user Token, scope string) bool {
	// scopes: user, moderator, admin
	// user > moderator > admin
	if user.Scope == "admin" {
		return true
	} else if user.Scope == "moderator" {
		if scope == "user" || scope == "moderator" {
			return true
		}
	} else if user.Scope == "user" {
		if scope == "user" {
			return true
		}
	}
	return user.Scope == scope
}
