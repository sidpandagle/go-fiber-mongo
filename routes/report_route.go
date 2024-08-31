package routes

import (
	"fibgo/controllers"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func ReportRoute(app *fiber.App) {

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWTSECRET"))},
	}))
	app.Post("/report", controllers.CreateReport)
	app.Get("/report/:reportId", controllers.GetAReport)
	app.Put("/report/:reportId", controllers.EditAReport)
	app.Delete("/report/:reportId", controllers.DeleteAReport)
	app.Get("/report", controllers.GetAllReports)
}
