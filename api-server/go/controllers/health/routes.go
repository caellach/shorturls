package health

import "github.com/gofiber/fiber/v2"

func getHealth(c *fiber.Ctx) error {
	return c.SendString("OK")
}
