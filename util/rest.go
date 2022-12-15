package util

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
)

func RestResponse(c *fiber.Ctx, status int, data interface{}) error {
	if data == nil && data == "" {
		data = []string{}
	}
	kind := reflect.TypeOf(data).Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		if reflect.ValueOf(data).Len() == 0 {
			data = []string{}
		}
	}
	success := status >= 200 && status < 300
	return c.Status(status).JSON(fiber.Map{
		"success": success,
		"data":    data,
	})
}
