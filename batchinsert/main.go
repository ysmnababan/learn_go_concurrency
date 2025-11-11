package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func insertConcurrent(db *gorm.DB, data []Employee, cnt int) {
	var wg sync.WaitGroup

	wg.Add(cnt)
	fmt.Println(len(data))
	var chs []chan Employee
	for range cnt {
		c := make(chan Employee, len(data)/cnt)
		chs = append(chs, c)
	}
	for i, e := range data {
		idx := i % cnt
		chs[idx] <- e
	}
	for i := range cnt {
		close(chs[i])
	}
	for i := range cnt {
		go func(ch <-chan Employee, dbs *gorm.DB) {
			defer wg.Done()
			session := dbs.Session(&gorm.Session{NewDB: true})
			affectedRows := 0
			for e := range ch {
				res := session.Create(&e)
				if res.Error != nil {
					log.Println(res.Error)
					continue
				}
				affectedRows++
			}
			fmt.Println("go routine no: ", i, affectedRows)
		}(chs[i], db)
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
