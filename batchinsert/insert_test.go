package main

import (
	"log"
	"testing"

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
			_ = db.Create(&d)
		}

	}

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

}
