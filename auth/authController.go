package auth

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"necross.it/backend/auth/user"
	"necross.it/backend/mail"
	"necross.it/backend/util"
	"strconv"
)

func (s *Service) LoginAuth(ctx *fiber.Ctx) error {
	var token = ctx.Get("Authorization")
	var authForm FormAuth
	_ = ctx.BodyParser(&authForm)
	user, token, err := s.Login(authForm.Email, []byte(authForm.Password))
	if err != nil {
		return util.RestResponse(ctx, 500, err.Error())
	}
	err = s.DeactivateTokens(user)
	auth, err := s.CreateToken(user, token)
	if err != nil {
		return util.RestResponse(ctx, 500, err.Error())
	}
	return ctx.Status(200).JSON(fiber.Map{
		"user":         user.ConvertToAuthJSON(),
		"access_token": auth.Token,
	})
}

func (s *Service) Logout(ctx *fiber.Ctx) error {
	var bearerToken = ctx.Get("Authorization")
	if bearerToken == "" || bearerToken == "Bearer" {
		return util.RestResponse(ctx, 401, "No token provided")
	}
	var token = bearerToken[7:]
	var authToken, err = s.GetByToken(token)
	if err != nil {
		return util.RestResponse(ctx, 401, "Invalid token")
	}
	user, err := s.User.GetUserById(strconv.Itoa(authToken.Id))
	if err != nil {
		return util.RestResponse(ctx, 401, "Could not find user")
	}
	err = s.DeactivateTokens(user)
	if err != nil {
		return util.RestResponse(ctx, 500, "Could not delete token")
	}
	return util.RestResponse(ctx, 200, "Erfolgreich abgemeldet!")
}

func (s *Service) UserMe(ctx *fiber.Ctx) error {
	var bearerToken = ctx.Get("Authorization")
	if bearerToken == "" || bearerToken == "Bearer" {
		return util.RestResponse(ctx, 401, "No token provided")
	}
	token := bearerToken[7:]
	authToken, err := s.GetByToken(token)
	isExpired := s.CheckTokenExpired(authToken.Token)
	if isExpired {
		return util.RestResponse(ctx, 401, "Token expired")
	}
	if err != nil {
		return util.RestResponse(ctx, 401, "Invalid token")
	}
	user, err := s.User.GetUserById(strconv.Itoa(authToken.UserId))
	if err != nil {
		return util.RestResponse(ctx, 401, "Invalid token")
	}
	return ctx.Status(200).JSON(user.ConvertToAuthJSON())
}

func (s *Service) RegisterController(ctx *fiber.Ctx) error {
	var registerForm RegisterForm
	_ = ctx.BodyParser(&registerForm)
	if s.UserExistsByEmail(registerForm.Email) {
		return util.RestResponse(ctx, 400, "Diese E-Mail Adresse ist bereits vergeben!")
	}
	if !util.IsValidEmail(registerForm.Email) {
		return util.RestResponse(ctx, 400, "Diese E-Mail Adresse ist ungültig!")
	}
	if len(registerForm.Password) < 8 {
		return util.RestResponse(ctx, 400, "Das Passwort muss mindestens 8 Zeichen lang sein!")
	}
	user, err := s.RegisterService(registerForm)
	if err != nil {
		return util.RestResponse(ctx, 500, "Es ist ein Fehler aufgetreten! Bitte versuche es später erneut!")
	}
	user, token, err := s.Login(registerForm.Email, []byte(registerForm.Password))
	auth, err := s.CreateToken(user, token)
	if err != nil {
		return util.RestResponse(ctx, 500, err.Error())
	}
	return ctx.Status(200).JSON(fiber.Map{
		"user":         user.ConvertToAuthJSON(),
		"access_token": auth.Token,
	})
}

func (s *Service) DeleteAccount(ctx *fiber.Ctx) error {
	var bearerToken = ctx.Get("Authorization")
	if bearerToken == "" || bearerToken == "Bearer" {
		return util.RestResponse(ctx, 401, "No token provided")
	}
	token := bearerToken[7:]
	authToken, err := s.GetByToken(token)
	if err != nil {
		return util.RestResponse(ctx, 401, "Invalid token")
	}
	user, err := s.User.GetUserById(strconv.Itoa(authToken.UserId))
	err = mail.SendDeleteAccountMail(user)
	if err != nil {
		return util.RestResponse(ctx, 500, "Could not send mail")
	}
	return util.RestResponse(ctx, 200, "Wir haben dir eine E-Mail zum Löschen deines Accounts geschickt!")
}

func (s *Service) AuthToken(ctx *fiber.Ctx) (Token, user.User, error) {
	var bearerToken = ctx.Get("Authorization")
	if bearerToken == "" || bearerToken == "Bearer" {
		return Token{}, user.User{}, errors.New("No token provided")
	}
	token := bearerToken[7:]
	authToken, err := s.GetByToken(token)
	emptyUser := user.User{}
	if err != nil {
		return Token{}, emptyUser, errors.New("Invalid token")
	}
	user, err := s.User.GetUserById(strconv.Itoa(authToken.UserId))
	if err != nil {
		return Token{}, emptyUser, errors.New("Invalid User")
	}
	return authToken, user, nil
}

func (s *Service) ChangeEmail(ctx *fiber.Ctx) error {
	_, user, err := s.AuthToken(ctx)
	var changeEmailForm ChangeEmailForm
	err = ctx.BodyParser(&changeEmailForm)
	if err != nil {
		return util.RestResponse(ctx, 400, "Invalid token")
	}
	err = s.User.ChangeEmail(changeEmailForm.Email, user)
	if err != nil {
		return util.RestResponse(ctx, 500, "Could not change email")
	}
	err = s.DeactivateTokens(user)
	if err != nil {
		return util.RestResponse(ctx, 500, "E-Mail konnte nicht geändert werden!")
	}
	return util.RestResponse(ctx, 200, "E-Mail wurde erfolgreich geändert, und alle anderen Geräte wurden abgemeldet!")
}

func (s *Service) ChangePassword(ctx *fiber.Ctx) error {
	_, user, err := s.AuthToken(ctx)
	var changePasswordForm ChangePasswordForm
	err = ctx.BodyParser(&changePasswordForm)
	if err != nil {
		return util.RestResponse(ctx, 400, "Invalid token")
	}
	if len(changePasswordForm.Password) < 8 {
		return util.RestResponse(ctx, 400, "Das Passwort muss mindestens 8 Zeichen lang sein!")
	}
	err = s.User.ChangePassword(changePasswordForm.Password, user)
	if err != nil {
		return util.RestResponse(ctx, 500, "Could not change password")
	}
	err = s.DeactivateTokens(user)
	if err != nil {
		return util.RestResponse(ctx, 500, "Sessions konnten nicht abgemeldet werden!")
	}
	return util.RestResponse(ctx, 200, "Passwort wurde erfolgreich geändert, und alle Geräte wurden abgemeldet!")
}

func (s *Service) ForgotPassword(ctx *fiber.Ctx) error {
	var forgotPasswordForm ChangeEmailForm
	err := ctx.BodyParser(&forgotPasswordForm)
	if err != nil {
		return util.RestResponse(ctx, 400, "Bitte gib eine E-Mail Adresse ein!")
	}
	user, err := s.GetUserByEmail(forgotPasswordForm.Email)
	token, err := s.User.ForgotPassword(user)
	if err != nil {
		return util.RestResponse(ctx, 500, "Es ist ein Fehler aufgetreten! Bitte versuche es später erneut!")
	}
	go mail.SendForgotPasswordMail(user, mail.ForgotPasswordMail{Link: "https://pricelesskeys.com/auth/password/reset/" + token})
	return util.RestResponse(ctx, 200, "Wir haben dir eine E-Mail zum Zurücksetzen deines Passworts geschickt!")
}

func (s *Service) ResetPassword(ctx *fiber.Ctx) error {
	var resetPasswordForm ResetPasswordForm
	err := ctx.BodyParser(&resetPasswordForm)
	if err != nil {
		return util.RestResponse(ctx, 400, "Bitte gib alle notwendigen Daten ein!")
	}
	user, err := s.GetUserByEmail(resetPasswordForm.Email)
	if !s.User.ValidatePasswordResetToken(resetPasswordForm.Token, user) {
		return util.RestResponse(ctx, 400, "Der Token ist ungültig")
	}
	err = s.User.ChangePassword(resetPasswordForm.Password, user)
	if err != nil {
		return util.RestResponse(ctx, 500, "Das Passwort konnte nicht zurückgesetzt werden!")
	}
	err = s.DeactivateTokens(user)
	err = s.User.DeleteForgotPasswordToken(user)
	if err != nil {
		return util.RestResponse(ctx, 500, "Das Passwort konnte nicht zurückgesetzt werden!")
	}
	return util.RestResponse(ctx, 200, "Das Passwort wurde erfolgreich zurückgesetzt!")
}

type ResetPasswordForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type ChangeEmailForm struct {
	Email string `json:"email"`
}

type ChangePasswordForm struct {
	Password string `json:"password"`
}

func (s *Service) DeleteAccountConfirm(ctx *fiber.Ctx) error {
	var bearerToken = ctx.Get("Authorization")
	if bearerToken == "" || bearerToken == "Bearer" {
		return util.RestResponse(ctx, 401, "Invalid token")
	}
	token := bearerToken[7:]
	authToken, err := s.GetByToken(token)
	if err != nil {
		return util.RestResponse(ctx, 401, "Invalid token")
	}
	err = s.User.DeleteUser(strconv.Itoa(authToken.UserId))
	if err != nil {
		return util.RestResponse(ctx, 500, "Could not delete user")
	}
	return util.RestResponse(ctx, 200, "Account deleted successfully!")

}
