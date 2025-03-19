package tables

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ExternalID uint       `gorm:"not null" json:"externalId"`
	Name       string     `json:"name"`
	Duration   int        `gorm:"not null" json:"duration"`
	Difficulty int        `gorm:"not null" json:"difficulty"`
	Provider   string     `gorm:"not null" json:"provider"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

func (Task) TableName() string {
	return "tb_tasks"
}

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	t.CreatedAt = &now
	t.UpdatedAt = &now
	return
}

func (t *Task) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	t.UpdatedAt = &now
	return
}
