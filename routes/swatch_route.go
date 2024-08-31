package routes

import (
	"fibgo/controllers"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func SwatchRoute(app *fiber.App) {

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWTSECRET"))},
	}))
	app.Post("/swatch-add", controllers.CreateSwatchList)
	app.Post("/swatch", controllers.CreateSwatch)
	app.Get("/swatch/:swatchId", controllers.GetASwatch)
	app.Put("/swatch/:swatchId", controllers.EditASwatch)
	app.Get("/swatch-like/:swatchId", controllers.IncrementLike)
	app.Delete("/swatch/:swatchId", controllers.DeleteASwatch)
	app.Get("/swatch", controllers.GetAllSwatch)
	app.Get("/swatch-filter", controllers.GetFilteredSwatch)
}
