package routes

import (
	"fibgo/controllers"

	"github.com/gofiber/fiber/v2"
)

func ReportRoute(app *fiber.App) {
	app.Post("/report", controllers.CreateReport)
	app.Get("/report/:reportId", controllers.GetAReport)
	app.Put("/report/:reportId", controllers.EditAReport)
	app.Delete("/report/:reportId", controllers.DeleteAReport)
	app.Get("/reports", controllers.GetAllReports)
}
