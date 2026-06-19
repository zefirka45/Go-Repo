package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Users struct {
	ID        uint16
	Userlogin string
	Username  string
	Role      string
	Password  string
	Status    bool
}

func ConnectPostgres() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("db_host"),
		os.Getenv("db_user"),
		os.Getenv("db_password"),
		os.Getenv("db_name"),
		os.Getenv("db_port"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db, err
}
func CreateTableUsers(db *gorm.DB) {
	db.Create(&Users{})

}
