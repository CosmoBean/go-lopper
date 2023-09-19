package model

import (
	"fmt"
	"os"

	"github.com/oklog/ulid/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Url struct {
	ID       ulid.ULID `json:"id" gorm:"primary_key" swaggertype:"string"`
	Redirect string    `json:"redirect"`
	Lopper   string    `json:"lopper" gorm:"unique;not null"`
	Clicked  uint64    `json:"clicked"`
	Random   bool      `json:"random"`
}

func Setup() {
	var (
		dbHost, dbName, dbUser, dbPassword, dbPort string
		ok                                         bool
	)

	if dbHost, ok = os.LookupEnv("DATABASE_HOST"); !ok {
		panic("DATABASE_HOST not found in env")
	}

	if dbName, ok = os.LookupEnv("DATABASE_NAME"); !ok {
		panic("DATABASE_NAME not found in env")
	}

	if dbUser, ok = os.LookupEnv("DATABASE_USER"); !ok {
		panic("DATABASE_USER not found in env")
	}

	if dbPassword, ok = os.LookupEnv("DATABASE_PASSWORD"); !ok {
		panic("DATABASE_PASSWORD not found in env")
	}

	if dbPort, ok = os.LookupEnv("DATABASE_PORT"); !ok {
		panic("DATABASE_PORT not found in env")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Url{})
	if err != nil {
		fmt.Println(err)
	}
}
