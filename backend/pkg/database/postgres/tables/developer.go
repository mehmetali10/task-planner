package tables

import (
	"time"

	"gorm.io/gorm"
)

type Developer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	Email     string    `gorm:"not null" json:"email"`
	Capacity  int       `gorm:"not null" json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Developer) TableName() string {
	return "tb_developers"
}

func (d *Developer) BeforeCreate(tx *gorm.DB) (err error) {
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	return
}

func (d *Developer) BeforeUpdate(tx *gorm.DB) (err error) {
	d.UpdatedAt = time.Now()
	return
}
