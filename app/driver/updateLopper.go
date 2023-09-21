package driver

import (
	"go-lopper/db"
	"go-lopper/model"
	"go-lopper/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

func UpdateLopper(ctx *fiber.Ctx, id ulid.ULID, updateRequest model.UrlRequest) (model.Url, error) {
	var url model.Url

	//lopper validations
	url, random, err := utils.ValidateLopper(updateRequest.Lopper)
	if url.ID.String() != id.String() {
		if err != nil {
			return url, err
		}

		url, err = db.GetUrl(id)
		if err != nil {
			return url, err
		}
	}
	url.Random = random

	url.Redirect = updateRequest.Redirect
	if url.Random {
		url.Lopper = utils.RandomUrl(8)
	} else {
		url.Lopper = updateRequest.Lopper
	}

	if err = db.UpdateUrl(url); err != nil {
		return url, err
	}

	return url, nil
}
