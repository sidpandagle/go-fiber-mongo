package routes

import (
	"fibgo/controllers"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func HabitRoute(app *fiber.App) {

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWTSECRET"))},
	}))
	app.Post("/habits", controllers.CreateHabits)
	app.Post("/habit", controllers.CreateHabit)
	app.Get("/habit/:habitId", controllers.GetAHabit)
	app.Put("/habit/:habitId", controllers.EditAHabit)
	app.Delete("/habit/:habitId", controllers.DeleteAHabit)
	app.Get("/habit", controllers.GetAllHabit)
}
