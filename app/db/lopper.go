package db

import (
	"go-lopper/model"

	"github.com/oklog/ulid/v2"
)

func GetAllUrls() ([]model.Url, error) {
	var urls []model.Url

	tx := db.Find(&urls)

	if tx.Error != nil {
		return []model.Url{}, tx.Error
	}

	return urls, nil
}

func GetUrl(id ulid.ULID) (model.Url, error) {
	var url model.Url

	tx := db.Where("id = ?", id).First(&url)

	if tx.Error != nil {
		return model.Url{}, tx.Error
	}

	return url, nil
}

func CreateUrl(url model.Url) error {
	tx := db.Create(&url)
	return tx.Error
}

func UpdateUrl(url model.Url) error {
	tx := db.Save(&url)
	return tx.Error
}

func DeleteUrlByLopper(lopper string) error {
	tx := db.Unscoped().Where("lopper = ?", lopper).Delete(&model.Url{})
	return tx.Error
}

func DeleteUrl(id ulid.ULID) error {
	tx := db.Unscoped().Where("id = ?", id).Delete(&model.Url{})
	return tx.Error
}

func FindUrlByLopper(Lopper string) (model.Url, bool, error) {
	var url model.Url
	tx := db.Where("lopper = ?", url).First(&url)
	if tx.Error != nil {
		return url, false, tx.Error
	}
	return url, true, tx.Error
}
