package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
func Connect(){
	dsn := "root:pass@tcp(127.0.0.1:3306)/institution?charset=utf8mb4&parseTime=True&loc=Local"
	db,err:=gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err!=nil{
		log.Fatal("cant connect to the database")
	}
	DB=db
	log.Print("Connected to the database sucessfully")

}