package server

import (
	"fmt"
	"go-lopper/model"
	"go-lopper/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

func getPing(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "pong"})
}

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

func getAllRedirects(ctx *fiber.Ctx) error {
	urls, err := model.GetAllUrls()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"message": "Something went wrong while fetching urls " + err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(urls)
}

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
