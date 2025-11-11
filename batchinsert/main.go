package main

import (
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func insertConcurrent(db *gorm.DB, data []Employee, cnt int) {
	var wg sync.WaitGroup
	wg.Add(cnt)

	channels := make([]chan Employee, cnt)
	for i := range cnt {
		channels[i] = make(chan Employee, len(data)/cnt)
	}

	// Distribute work
	for i, e := range data {
		idx := i % cnt
		channels[idx] <- e
	}

	for i := range cnt {
		close(channels[i])
	}

	// Start goroutines
	for i := range cnt {
		ch := channels[i] // capture i
		go func(ch <-chan Employee) {
			defer wg.Done()
			session := db.Session(&gorm.Session{NewDB: true}) // create new db session for each go routine
			for e := range ch {
				if err := session.Create(&e).Error; err != nil {
					log.Println(err)
				}
			}
		}(ch)
	}

	wg.Wait()
}
func main() {

	dsn := "host=localhost user=postgres password=postgres dbname=malut port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	err = db.AutoMigrate(&Employee{})

	if err != nil {
		panic(err)
	}
	db.Where("id is not null").Delete(&Employee{})
}
