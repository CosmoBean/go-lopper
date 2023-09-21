package utils

import (
	"fmt"
	"go-lopper/db"
	"go-lopper/model"
)

func ValidateLopper(lopper string) (url model.Url, random bool, err error) {
	lenLopper := len(lopper)
	if lenLopper == 0 {
		return url, true, nil
	} else {
		// length validation
		if lenLopper < 4 {
			return url, false, fmt.Errorf("lopper length should be at least 4 characters")
		}
		//unique validation
		url, ok, err := db.GetUrlByLopper(lopper)
		if err == nil && ok {
			return url, false, fmt.Errorf("lopper already exists")
		}
	}
	return url, false, nil
}
