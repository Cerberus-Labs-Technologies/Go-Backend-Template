package user

import (
	"golang.org/x/crypto/bcrypt"
	"necross.it/backend/database"
	"necross.it/backend/util"
)

type Service struct {
	Server database.Server
}

func (s *Service) GetUserById(id string) (User, error) {
	var user User
	err := s.Server.DB.Get(&user, "SELECT * FROM users WHERE id = ?", id)
	return user, err
}

func (s *Service) Register(u User) error {
	_, err := s.Server.DB.NamedExec("INSERT INTO users (`name`, email, password) VALUES (:name, :email, :password)", u)
	return err
}

func (s *Service) ForgotPassword(u User) (string, error) {
	token := util.GenerateRandomTokenWithLength(32)
	forgotForm := ForgotPassword{
		UserID:     u.Id,
		ResetToken: token,
	}
	err := s.Create(forgotForm)
	return token, err
}

func (s *Service) ChangePassword(newPassword string, u User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	_, err = s.Server.DB.Exec("UPDATE users SET password = ? WHERE id = ?", hashedPassword, u.Id)
	return err
}

func (s *Service) Create(f ForgotPassword) error {
	_, err := s.Server.DB.NamedExec("INSERT INTO forgot_password (userId, resetToken) VALUES (:userId, :resetToken)", f)
	return err
}

func (s *Service) ValidatePasswordResetToken(token string, u User) bool {
	var forgotPassword ForgotPassword
	err := s.Server.DB.Get(&forgotPassword, "SELECT * FROM forgot_password WHERE userId = ? AND resetToken = ?", u.Id, token)
	return err == nil
}

func (s *Service) DeleteForgotPasswordToken(u User) error {
	_, err := s.Server.DB.Exec("DELETE FROM forgot_password WHERE userId = ?", u.Id)
	return err
}

func (s *Service) ChangeEmail(email string, u User) error {
	_, err := s.Server.DB.Exec("UPDATE users SET email = ? WHERE id = ?", email, u.Id)
	return err
}

func (s *Service) GetAllUsers() (int, error) {
	var users []User
	err := s.Server.DB.Select(&users, "SELECT * FROM users")
	return len(users), err
}

func (s *Service) DeleteUser(id string) error {
	_, err := s.Server.DB.Exec("DELETE FROM users WHERE id = ?", id)
	_, err = s.Server.DB.Exec("DELETE FROM auth_tokens WHERE user_id = ?", id)
	_, err = s.Server.DB.Exec("DELETE FROM verify_mail WHERE userId = ?", id)
	_, err = s.Server.DB.Exec("DELETE FROM watchers WHERE userID = ?", id)
	return err
}
