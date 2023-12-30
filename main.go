package main

import (
    "fibgo/configs"
    "fibgo/routes"
	"github.com/gofiber/fiber/v2"
)

func hello(c *fiber.Ctx) error{
	return c.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB"})
}

func main(){
	app := fiber.New()

    configs.ConnectDB()

    routes.UserRoute(app)

	app.Get("/", hello)

	// app.Static("/", "./public") 

	app.Listen(":3000")
}