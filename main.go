package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sidpandagle/fiber/sample"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	return port
}

func main() {
	app := fiber.New()

	run();

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, Railway!",
		})
	})

	app.Listen(getPort())
}
