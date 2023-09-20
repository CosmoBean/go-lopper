package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"

	_ "go-lopper/docs"
)

func SetupAndServe() {

	router := fiber.New()

	//cors settings
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	//docs
	router.Get("/docs/*", swagger.HandlerDefault) // default

	//routes
	router.Get("/health", getPing)
	router.Get("/lopper", getAllRedirects)
	router.Get("/lopper/:id", getRedirectUrl)
	router.Post("/lopper", createRedirectUrl)
	router.Patch("/lopper", updateRedirectUrl)
	router.Delete("/lopper/:id", deleteRedirectUrl)
	router.Delete("/lopper", deleteRedirectUrlByLopper)

	router.Get("r/:redirect", redirect)

	router.Listen(":9020")
}
