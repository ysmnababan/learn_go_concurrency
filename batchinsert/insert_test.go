package main

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func BenchmarkInsertSingle(b *testing.B) {
	datas := []Employee{}
	n := 10000
	for range n {
		emp := Employee{
			NIP:          uuid.NewString(),
			Name:         "some name",
			OrgUnit:      "org unit",
			EmployeeType: "employee type",
			Status:       "status",
			IsActive:     true,
		}
		datas = append(datas, emp)
	}
	dsn := "host=localhost user=postgres password=postgres dbname=malut port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		db.Where("id is not null").Delete(&Employee{})
		b.StartTimer()
		for _, d := range datas {
			tx := db.Create(&d)
			if tx.Error != nil {
				log.Println(tx.Error)
			}
		}

	}
	db.Where("id is not null").Delete(&Employee{})

}

func BenchmarkInsertBatch(b *testing.B) {

	datas := []Employee{}
	n := 10000
	for range n {
		emp := Employee{
			NIP:          uuid.NewString(),
			Name:         "some name",
			OrgUnit:      "org unit",
			EmployeeType: "employee type",
			Status:       "status",
			IsActive:     true,
		}
		datas = append(datas, emp)
	}
	dsn := "host=localhost user=postgres password=postgres dbname=malut port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		db.Where("id is not null").Delete(&Employee{})
		b.StartTimer()
		_ = db.Create(&datas)

	}
	db.Where("id is not null").Delete(&Employee{})

}
func BenchmarkInsertGoroutine(b *testing.B) {
	datas := []Employee{}
	n := 10000
	for range n {
		emp := Employee{
			NIP:          uuid.NewString(),
			Name:         "some name",
			OrgUnit:      "org unit",
			EmployeeType: "employee type",
			Status:       "status",
			IsActive:     true,
		}
		datas = append(datas, emp)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		dsn := "host=localhost user=postgres password=postgres dbname=malut port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println(err)
			return
		}
		sqlDB, _ := db.DB()
		gcnt := 5
		sqlDB.SetMaxIdleConns(gcnt)
		sqlDB.SetMaxOpenConns(gcnt)
		sqlDB.SetConnMaxLifetime(time.Hour)
		db.Where("id is not null").Delete(&Employee{})
		b.StartTimer()

		insertConcurrent(db, datas, gcnt)
		sqlDB.Close()
	}
	// db.Where("id is not null").Delete(&Employee{})

}
