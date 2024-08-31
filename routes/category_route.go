package routes

import (
	"fibgo/controllers"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func CategoryRoute(app *fiber.App) {

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWTSECRET"))},
	}))
	app.Post("/category", controllers.CreateCategory)
	app.Get("/category/:categoryId", controllers.GetACategory)
	app.Put("/category/:categoryId", controllers.EditACategory)
	app.Delete("/category/:categoryId", controllers.DeleteACategory)
	app.Get("/category", controllers.GetAllCategory)
}
