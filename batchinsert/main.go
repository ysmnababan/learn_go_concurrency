package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func singleInsert(db *gorm.DB, data []*Employee) {
	for _, d := range data {
		db.Create(d)
	}
}
func main() {

	dsn := "host=localhost user=postgres password=postgres dbname=malut port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&Employee{})
	if err != nil {
		panic(err)
	}
	db.Where("id is not null").Delete(&Employee{})
}
