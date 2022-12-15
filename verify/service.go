package verify

import (
	"github.com/gofiber/fiber/v2"
	"necross.it/backend/auth/user"
	"necross.it/backend/database"
	"necross.it/backend/mail"
	"necross.it/backend/util"
	"strconv"
)

type Service struct {
	Server database.Server
}

func (s *Service) RegisterRoutes(verifyGroup fiber.Router) {
	verifyGroup.Post("/:verifyToken", s.VerifyUser)
}

func (s *Service) verifyEmail(email string) error {
	verifyMail, err := s.getVerifyTokenByEmail(email)
	_, err = s.Server.DB.Exec("UPDATE users SET email_verified_at = NOW() WHERE email = ?", email)
	if err != nil {
		return err
	}
	_, err = s.Server.DB.Exec("DELETE FROM verify_mail WHERE verifyToken = ?", verifyMail.VerifyToken)
	return err
}

// VerifyTokenExists func verifyToken exists
func (s *Service) verifyTokenExists(token string) (bool, error) {
	var count int
	err := s.Server.DB.Get(&count, "SELECT COUNT(*) FROM verify_mail WHERE verifyToken = ?", token)
	return count > 0, err
}

func (s *Service) verifyEmailByToken(token string) error {
	user, err := s.getUserByToken(token)
	_, err = s.Server.DB.Exec("UPDATE users SET email_verified_at = NOW() WHERE id = ?", strconv.Itoa(user.Id))
	if err != nil {
		return err
	}
	_, err = s.Server.DB.Exec("DELETE FROM verify_mail WHERE verifyToken = ?", token)
	return err
}

func (s *Service) getUserByToken(token string) (user.User, error) {
	var verifyModel mail.MailVerify
	var user user.User
	err := s.Server.DB.Get(&verifyModel, "SELECT * FROM verify_mail WHERE verifyToken = ?", token)
	if err != nil {
		return user, err
	}
	err = s.Server.DB.Get(&user, "SELECT * FROM users WHERE id = ?", strconv.Itoa(verifyModel.UserID))
	return user, err
}

func (s *Service) getVerifyTokenByEmail(email string) (mail.MailVerify, error) {
	// get verify token from sql
	var verifyModel mail.MailVerify
	err := s.Server.DB.Get(&verifyModel, "SELECT * FROM verify_mail WHERE email = ?", email)
	return verifyModel, err
}

func (s *Service) createVerifyEmail(user user.User) (mail.MailVerify, error) {
	verifyMail := mail.MailVerify{
		UserID:      user.Id,
		VerifyToken: util.GenerateRandomTokenWithLength(32),
	}
	_, err := s.Server.DB.NamedExec("INSERT INTO verify_mail (userId, verifyToken) VALUES (:userId, :verifyToken)", verifyMail)
	return verifyMail, err
}

func (s *Service) verify(u user.User) error {
	return s.verifyEmail(u.Email)
}

func (s *Service) SendVerifyEmail(u user.User) error {
	verify, err := s.createVerifyEmail(u)
	if err != nil {
		return err
	}
	err = mail.SendVerifyMail(u, verify)
	return err
}
