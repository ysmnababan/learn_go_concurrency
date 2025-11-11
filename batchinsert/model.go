package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	NIP          string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"nip"`
	Name         string     `gorm:"type:varchar(255);not null" json:"name"`
	OrgUnit      string     `gorm:"type:varchar(255);not null;index" json:"org_unit"`
	EmployeeType string     `gorm:"type:text;not null" json:"employee_type"` // PNS, PPPK, CONTRACT, OUTSOURCED
	Position     *string    `gorm:"type:varchar(255)" json:"position,omitempty"`
	GolonganPNS  *string    `gorm:"type:varchar(50)" json:"golongan_pns,omitempty"`
	GolonganPPPK *string    `gorm:"type:varchar(50)" json:"golongan_pppk,omitempty"`
	PlaceOfBirth *string    `gorm:"type:varchar(100)" json:"place_of_birth,omitempty"`
	DateOfBirth  *time.Time `gorm:"type:date" json:"date_of_birth,omitempty"`
	Status       string     `gorm:"type:text;not null;index" json:"status"` // ACTIVE, ON_LEAVE, ...
	Salary       float64    `gorm:"type:numeric(12,2);default:0.00" json:"salary"`
	IsActive     bool       `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate hook to set UUID if not set (optional)
func (e *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}
func (Employee) TableName() string {
	return "temp_employee"
}
