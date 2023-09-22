package server

import (
	"fmt"
	"go-lopper/db"
	"go-lopper/driver"
	"go-lopper/model"
	"go-lopper/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

// GetPing godoc
// @Summary Check API health
// @Description Returns a pong response if API is healthy
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "message:alive"
// @Router /health [get]
func getPing(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "alive"})
}

// Redirect godoc
// @Summary Redirect to original URL
// @Description Redirects the user to the original URL based on the lopper value
// @Tags redirect
// @Param redirect path string true "lopper value"
// @Produce  json
// @Success 307 {string} string "Redirected"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /r/{redirect} [get]
func redirect(ctx *fiber.Ctx) error {
	lopper := ctx.Params("redirect")
	redirectUrl, _, err := db.GetUrlByLopper(lopper)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Error while finding by URL " + err.Error()})
	}

	redirectUrl.Clicked += 1
	if err := db.UpdateUrl(redirectUrl); err != nil {
		fmt.Println("failed to save clicked to db ", err)
	}

	return ctx.Redirect(redirectUrl.Redirect, fiber.StatusTemporaryRedirect)

}

// GetAllRedirects godoc
// @Summary Retrieve all redirects
// @Description Get a list of all redirect URLs
// @Tags redirect
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Url
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /lopper [get]
func getAllRedirects(ctx *fiber.Ctx) error {
	urls, err := db.GetAllUrls()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Something went wrong while fetching urls " + err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(urls)
}

// GetRedirectUrl godoc
// @Summary Get redirect URL by ID
// @Description Retrieve a specific redirect URL by its ID
// @Tags redirect
// @Param id path string true "ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Url
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /lopper/{id} [get]
func getRedirectUrl(ctx *fiber.Ctx) error {
	id, err := ulid.Parse(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Something went wrong while fetching redirect url " + err.Error()})

	}

	redirectUrl, err := db.GetUrl(id)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Something went wrong while fetching redirect url " + err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(redirectUrl)
}

// CreateRedirectUrl godoc
// @Summary Create a new redirect URL
// @Description Create a new redirect URL with optional custom lopper
// @Tags redirect
// @Accept  json
// @Produce  json
// @Param url body model.UrlRequest true "request body"
// @Success 201 {object} model.Url
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /lopper [post]
func createRedirectUrl(ctx *fiber.Ctx) error {
	var url model.Url
	if err := ctx.BodyParser(&url); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Something went wrong while parsing body " + err.Error()})
	}

	if ok := utils.ValidateUrlString(url.Redirect); !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Invalid redirect url"})
	}

	//lopper validations
	if _, random, err := utils.ValidateLopper(url.Lopper); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Invalid lopper " + err.Error()})
	} else {
		url.Random = random
	}

	url.ID = ulid.Make()
	if url.Random {
		url.Lopper = utils.RandomUrl(8)
	}

	if err := db.CreateUrl(url); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Something went wrong while creating shortened url " + err.Error()})

	}

	return ctx.Status(fiber.StatusCreated).JSON(url)
}

// UpdateRedirectUrl godoc
// @Summary Update a redirect URL by ID
// @Description Update an existing redirect URL by its ID
// @Tags redirect
// @Param id path string true "ID"
// @Accept  json
// @Produce  json
// @Param url body model.UrlRequest true "request body"
// @Success 200 {object} model.Url
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /lopper/{id} [put]
func updateRedirectUrl(ctx *fiber.Ctx) error {
	var urlRequest model.UrlRequest
	id, err := ulid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Please pass the ID parameter properly  " + err.Error()})
	}

	if err := ctx.BodyParser(&urlRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Something went wrong while parsing body " + err.Error()})
	}

	if ok := utils.ValidateUrlString(urlRequest.Redirect); !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Invalid redirect url"})
	}

	url, err := driver.UpdateLopper(ctx, id, urlRequest)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Something went wrong while updating shortened url " + err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(url)

}

// DeleteRedirectUrl godoc
// @Summary Delete a redirect URL by ID
// @Description Delete a specific redirect URL by its ID
// @Tags redirect
// @Param id path string true "ID"
// @Produce  json
// @Success 204 "Successfully Deleted"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /lopper/{id} [delete]
func deleteRedirectUrl(ctx *fiber.Ctx) error {
	id, err := ulid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Please pass the ID parameter properly " + err.Error()})
	}

	if err := db.DeleteUrl(id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Something went wrong while deleting shortened url " + err.Error()})
	}

	return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Successfully deleted"})

}

// DeleteRedirectUrlByLopper godoc
// @Summary Delete a redirect URL by lopper
// @Description Delete a specific redirect URL by its lopper value
// @Tags redirect
// @Param lopper query string true "lopper value"
// @Produce  json
// @Success 204 "Successfully Deleted"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /lopper [delete]
func deleteRedirectUrlByLopper(ctx *fiber.Ctx) error {
	lopper := ctx.Query("lopper")
	if len(lopper) < 4 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Lopper should be at least 4 characters long"})
	}

	if err := db.DeleteUrlByLopper(lopper); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Something went wrong while deleting shortened url " + err.Error()})
	} else {
		return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Successfully deleted"})
	}

}
