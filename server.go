package main

import (
	"leisurely/database"
	"leisurely/database/models"
	"leisurely/modules/plan"
	"leisurely/modules/event"
	"leisurely/modules/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoute(app fiber.Router) {
	user.UserRouter(app)
	plan.PlanRouter(app)
	event.EventRouter(app)
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	v1 := app.Group("/v1")
	setupRoute(v1)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	database.DB.AutoMigrate(&models.User{}, &models.Tags{}, &models.Event{}, &models.Transportation{},
		&models.Plan{}, &models.Preference{}, &models.EventTag{})

	app.Listen(":3000")
}
