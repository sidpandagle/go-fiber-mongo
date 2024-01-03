package main

import (
	"fibgo/configs"
	"fibgo/routes"
	"os"

	"github.com/gofiber/fiber/v2"
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

func hello(c *fiber.Ctx) error {
	return c.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB"})
}

func main() {
	app := fiber.New()

	configs.ConnectDB()

	routes.UserRoute(app)
	routes.ReportRoute(app)

	app.Get("/", hello)

	// app.Static("/", "./public")

	app.Listen(getPort())
}
