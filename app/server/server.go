package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"

	_ "go-lopper/docs"
	"go-lopper/utils"
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

	//API routes
	api := router.Group("/api/v1")
	api.Get("/health", getPing)
	api.Get("/loppers", getAllRedirects)
	api.Get("/loppers/:id", getRedirectUrl)
	api.Post("/loppers", createRedirectUrl)
	api.Put("/loppers/:id", updateRedirectUrl)
	api.Delete("/loppers/:id", deleteRedirectUrl)
	api.Delete("/loppers", deleteRedirectUrlByLopper)

	router.Get("r/:redirect", redirect)

	router.Listen(utils.GetEnvWithDefault("FIBER_PORT", ":9020"))
}
