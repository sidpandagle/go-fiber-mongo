package main

import (
	"fibgo/configs"
	"fibgo/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3001"
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

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Use(
		logger.New(), // add Logger middleware
	)

	configs.ConnectDB()

	routes.HabitRoute(app)
	routes.SwatchRoute(app)
	routes.CategoryRoute(app)
	routes.ReportRoute(app)
	routes.UserRoute(app)

	app.Get("/", hello)

	// app.Static("/", "./public")

	app.Listen(getPort())
}
