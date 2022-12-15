package verify

import (
	"github.com/gofiber/fiber/v2"
	"necross.it/backend/util"
)

func (s *Service) VerifyUser(ctx *fiber.Ctx) error {
	verifyToken := ctx.Params("verifyToken")
	exists, err := s.verifyTokenExists(verifyToken)
	if !exists {
		return util.RestResponse(ctx, 404, "This verification link is invalid!")
	}
	err = s.verifyEmailByToken(verifyToken)
	if err != nil {
		return util.RestResponse(ctx, 500, "An error occurred while verifying your email!")
	}
	return util.RestResponse(ctx, 200, "Your email has been verified successfully!")
}
