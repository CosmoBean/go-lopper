package model

import (
	"github.com/oklog/ulid/v2"
)

type Url struct {
	ID       ulid.ULID `json:"id" gorm:"primary_key" swaggertype:"string"`
	Redirect string    `json:"redirect"`
	Lopper   string    `json:"lopper" gorm:"unique;not null"`
	Clicked  uint64    `json:"clicked" swaggerignore:"true"`
	Random   bool      `json:"random"`
}
