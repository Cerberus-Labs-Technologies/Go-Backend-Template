package auth

import (
	"necross.it/backend/auth/user"
	"strconv"
	"time"
)

func (s *Service) CreateToken(user user.User, token string) (Token, error) {

	arg := map[string]interface{}{
		"user_id": user.Id,
		"scope":   user.Scope,
		"token":   token,
		"active":  true,
	}
	_, err := s.Server.DB.NamedExec(
		"INSERT INTO auth_tokens (user_id, scope, token, active) VALUES (:user_id, :scope, :token, :active)",
		arg)
	var authToken Token
	err = s.Server.DB.Get(&authToken, "SELECT * FROM auth_tokens WHERE token = ? AND user_id = ?", token, strconv.Itoa(user.Id))
	return authToken, err
}

func (s *Service) CheckTokenExpired(token string) bool {
	var authToken Token
	_ = s.Server.DB.Get(&authToken, "SELECT * FROM auth_tokens WHERE token = ?", token)

	// check if expiresAt has been reached or not
	if authToken.ExpiresAt.Unix() < time.Now().Unix() {
		_ = s.DeactivateToken(token)
		return true // is expired
	}

	return false
}

func (s *Service) DeactivateTokens(user user.User) error {
	_, err := s.Server.DB.Exec("UPDATE auth_tokens SET active = 0 WHERE user_id = ?", strconv.Itoa(user.Id))
	return err
}

func (s *Service) DeactivateToken(token string) error {
	_, err := s.Server.DB.Exec("UPDATE auth_tokens SET active = 0 WHERE token = ?", token)
	return err
}

func (s *Service) GetByToken(token string) (Token, error) {
	var authToken Token
	err := s.Server.DB.Get(&authToken, "SELECT * FROM auth_tokens WHERE token = ?", token)
	return authToken, err
}
