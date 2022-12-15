package verify

import (
	"github.com/gofiber/fiber/v2"
	"necross.it/backend/util"
)

func (s *Service) VerifyUser(ctx *fiber.Ctx) error {
	verifyToken := ctx.Params("verifyToken")
	exists, err := s.verifyTokenExists(verifyToken)
	if !exists {
		return util.RestResponse(ctx, 404, "Dein Verifizierungslink ist ungültig")
	}
	err = s.verifyEmailByToken(verifyToken)
	if err != nil {
		return util.RestResponse(ctx, 500, "Es ist ein Fehler bei der Verifizierung aufgetreten")
	}
	return util.RestResponse(ctx, 200, "Deine E-Mail wurde erfolgreich verifiziert")
}
