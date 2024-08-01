package routes

import (
	"fibgo/controllers"

	"github.com/gofiber/fiber/v2"
)

func HabitRoute(app *fiber.App) {
	app.Post("/habit", controllers.CreateHabit)
	app.Get("/habit/:habitId", controllers.GetAHabit)
	app.Put("/habit/:habitId", controllers.EditAHabit)
	app.Delete("/habit/:habitId", controllers.DeleteAHabit)
	app.Get("/habit", controllers.GetAllHabit)
}
