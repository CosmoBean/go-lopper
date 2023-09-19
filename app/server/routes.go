package server

import (
	"fmt"
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
// @Success 200 {object} map[string]interface{} "message:pong"
// @Router /ping [get]
func getPing(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "pong"})
}

// Redirect godoc
// @Summary Redirect to original URL
// @Description Redirects the user to the original URL based on the lopper value
// @Tags redirect
// @Param redirect path string true "Lopper Value"
// @Produce  json
// @Success 307 {string} string "Redirected"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /redirect/{redirect} [get]
func redirect(ctx *fiber.Ctx) error {
	lopper := ctx.Params("redirect")
	redirectUrl, _, err := model.FindUrlByLopper(lopper)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Error while finding by URL " + err.Error()})
	}

	redirectUrl.Clicked += 1
	if err := model.UpdateUrl(redirectUrl); err != nil {
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
// @Router /redirects [get]
func getAllRedirects(ctx *fiber.Ctx) error {
	urls, err := model.GetAllUrls()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Something went wrong while fetching urls " + err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(urls)
}

// GetRedirectUrl godoc
// @Summary Get specific redirect URL
// @Description Retrieve a specific redirect URL by its ID
// @Tags redirect
// @Param id path string true "URL ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Url
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /redirect/{id} [get]
func getRedirectUrl(ctx *fiber.Ctx) error {
	id, err := ulid.Parse(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Something went wrong while fetching redirect url " + err.Error()})

	}

	redirectUrl, err := model.GetUrl(id)

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
// @Param url body model.Url true "URL Model"
// @Success 201 {object} model.Url
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /redirect [post]
func createRedirectUrl(ctx *fiber.Ctx) error {
	var url model.Url
	if err := ctx.BodyParser(&url); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Something went wrong while parsing body " + err.Error()})
	}

	//lopper validations
	lenLopper := len(url.Lopper)
	if lenLopper > 0 {
		url.Random = false
		if len(url.Lopper) < 4 {
			return ctx.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{"message": "Lopper should be at least 4 characters long"})
		}
		existingUrl, ok, err := model.FindUrlByLopper(url.Lopper)
		if err == nil && ok {
			return ctx.Status(fiber.StatusConflict).JSON(existingUrl)
		}
	}

	url.ID = ulid.Make()
	if url.Random {
		url.Lopper = utils.RandomUrl(8)
	}

	if err := model.CreateUrl(url); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Something went wrong while creating shortened url " + err.Error()})

	}

	return ctx.Status(fiber.StatusCreated).JSON(url)
}

// UpdateRedirectUrl godoc
// @Summary Update a redirect URL
// @Description Update an existing redirect URL by its model
// @Tags redirect
// @Accept  json
// @Produce  json
// @Param url body model.Url true "URL Model"
// @Success 200 {object} model.Url
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /redirect [put]
func updateRedirectUrl(ctx *fiber.Ctx) error {
	var url model.Url
	if err := ctx.BodyParser(&url); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Something went wrong while parsing body " + err.Error()})
	}

	if err := model.UpdateUrl(url); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Something went wrong while updating shortened url " + err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(url)

}

// DeleteRedirectUrl godoc
// @Summary Delete a redirect URL by ID
// @Description Delete a specific redirect URL by its ID
// @Tags redirect
// @Param id path string true "URL ID"
// @Produce  json
// @Success 204 {object} map[string]interface{} "Successfully Deleted"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /redirect/{id} [delete]
func deleteRedirectUrl(ctx *fiber.Ctx) error {
	id, err := ulid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Something went wrong while parsing body " + err.Error()})
	}

	if err := model.DeleteUrl(id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Something went wrong while deleting shortened url " + err.Error()})
	}

	return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Succesfully deleted"})

}

// DeleteRedirectUrlByLopper godoc
// @Summary Delete a redirect URL by lopper
// @Description Delete a specific redirect URL by its lopper value
// @Tags redirect
// @Param lopper query string true "Lopper Value"
// @Produce  json
// @Success 204 {object} map[string]interface{} "Successfully Deleted"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /redirect [delete]
func deleteRedirectUrlByLopper(ctx *fiber.Ctx) error {
	lopper := ctx.Query("lopper")
	if len(lopper) < 4 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"message": "Lopper should be at least 4 characters long"})
	}

	if err := model.DeleteUrlByLopper(lopper); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Something went wrong while deleting shortened url " + err.Error()})
	} else {
		return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Succesfully deleted"})
	}

}
