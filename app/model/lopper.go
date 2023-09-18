package model

import "github.com/oklog/ulid/v2"

func GetAllUrls() ([]Url, error) {
	var urls []Url

	tx := db.Find(&urls)

	if tx.Error != nil {
		return []Url{}, tx.Error
	}

	return urls, nil
}

func GetUrl(id ulid.ULID) (Url, error) {
	var url Url

	tx := db.Where("id = ?", id).First(&url)

	if tx.Error != nil {
		return Url{}, tx.Error
	}

	return url, nil
}

func CreateUrl(url Url) error {
	tx := db.Create(&url)
	return tx.Error
}

func UpdateUrl(url Url) error {
	tx := db.Save(&url)
	return tx.Error
}

func DeleteUrlByLopper(lopper string) error {
	tx := db.Unscoped().Where("lopper = ?", lopper).Delete(&Url{})
	return tx.Error
}

func DeleteUrl(id ulid.ULID) error {
	tx := db.Unscoped().Where("id = ?", id).Delete(&Url{})
	return tx.Error
}

func FindUrlByLopper(Lopper string) (Url, bool, error) {
	var url Url
	tx := db.Where("lopper = ?", url).First(&url)
	if tx.Error != nil {
		return url, false, tx.Error
	}
	return url, true, tx.Error
}
