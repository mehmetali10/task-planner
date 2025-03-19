package tables

import (
	"time"

	"gorm.io/gorm"
)

type Assignment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TaskID      uint      `gorm:"not null" json:"task_id"`
	DeveloperID uint      `gorm:"not null" json:"developer_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Assignment) TableName() string {
	return "tb_assignments"
}

func (a *Assignment) BeforeCreate(tx *gorm.DB) (err error) {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return
}

func (a *Assignment) BeforeUpdate(tx *gorm.DB) (err error) {
	a.UpdatedAt = time.Now()
	return
}
