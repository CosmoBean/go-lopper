package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupAndServe() {

	router := fiber.New()

	//cors settings
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	//routes
	router.Get("/", getPing)
	router.Get("/lopper", getAllRedirects)
	router.Get("/lopper/:id", getRedirectUrl)
	router.Post("/lopper", createRedirectUrl)
	router.Patch("/lopper", updateRedirectUrl)
	router.Delete("/lopper/:id", deleteRedirectUrl)

	router.Get("r/:redirect", redirect)

	router.Listen(":9020")
}
