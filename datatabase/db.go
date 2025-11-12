package datatabase

import (
	"fmt"
	models "jwt-authentication/Models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbInit() {
	var err error

	dsn :="host=localhost user=postgres password=trees5000 dbname=task"

	DB,err=gorm.Open(postgres.Open(dsn),&gorm.Config{})

	if err!=nil{
		fmt.Println("Cannot connect to datatabase")
		return
	}

	fmt.Println("Successfully connected to database")

	DB.AutoMigrate(&models.User{})

}