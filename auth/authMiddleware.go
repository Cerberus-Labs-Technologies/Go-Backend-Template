package auth

import (
	"github.com/gofiber/fiber/v2"
	"necross.it/backend/auth/user"
	"necross.it/backend/util"
	"strconv"
)

func (s *Service) AccessWithGuard(ctx *fiber.Ctx, scope string) error {
	var bearerToken = ctx.Get("Authorization")
	if bearerToken == "" {
		return util.RestResponse(ctx, 401, "Unauthorized")
	}
	token := bearerToken[7:]
	authToken, err := s.GetByToken(token)
	if err != nil {
		return util.RestResponse(ctx, 401, "Unauthorized")
	}
	isExpired := s.CheckTokenExpired(authToken.Token)
	if isExpired {
		return util.RestResponse(ctx, 401, "Unauthorized")
	}
	hasAccess := s.CheckScope(authToken, scope)
	if !hasAccess {
		return util.RestResponse(ctx, 401, "Ungenügende Berechtigung!")
	}
	return ctx.Next()
}

func (s *Service) AccessWithPermission(ctx *fiber.Ctx, permission string) error {
	var bearerToken = ctx.Get("Authorization")
	if bearerToken == "" {
		return util.RestResponse(ctx, 401, "Unauthorized")
	}
	token := bearerToken[7:]
	authToken, err := s.GetByToken(token)
	if err != nil {
		return util.RestResponse(ctx, 401, "Unauthorized")
	}
	isExpired := s.CheckTokenExpired(authToken.Token)
	if isExpired {
		return util.RestResponse(ctx, 401, "Unauthorized")
	}

	hasAccess := s.HasPermission(authToken, permission)
	if !hasAccess {
		return util.RestResponse(ctx, 401, "Unauthorized")
	}
	return ctx.Next()
}

func (s *Service) IsLoggedIn(ctx *fiber.Ctx) error {
	var bearerToken = ctx.Get("Authorization")
	if bearerToken == "" {
		return util.RestResponse(ctx, 401, "Unauthorized")
	}
	token := bearerToken[7:]
	authToken, err := s.GetByToken(token)
	if err != nil {
		return util.RestResponse(ctx, 401, "Unauthorized")
	}
	isExpired := s.CheckTokenExpired(authToken.Token)
	if isExpired {
		return util.RestResponse(ctx, 401, "Unauthorized")
	}
	return ctx.Next()
}

func (s *Service) GetUserTokenFromSession(ctx *fiber.Ctx) (user.User, Token, error) {
	var bearerToken = ctx.Get("Authorization")
	token := bearerToken[7:]
	authToken, err := s.GetByToken(token)
	user, err := s.User.GetUserById(strconv.Itoa(authToken.UserId))
	return user, authToken, err
}

func (s *Service) IsAdmin(ctx *fiber.Ctx) error {
	return s.AccessWithGuard(ctx, "admin")
}
