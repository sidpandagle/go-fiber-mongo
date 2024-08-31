package routes

import (
	"fibgo/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app *fiber.App) {

	app.Get("/login", controllers.Login)
	app.Get("/restricted", controllers.Restricted)
}
