package tables

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ExternalID string    `gorm:"not null;unique" json:"external_id"`
	Name       string    `gorm:"not null" json:"name"`
	Duration   int       `gorm:"not null" json:"duration"`
	Difficulty int       `gorm:"not null" json:"difficulty"`
	Provider   string    `gorm:"not null" json:"provider"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Task) TableName() string {
	return "tb_tasks"
}

// BeforeCreate hook: Task oluşturulmadan önce çalışır.
func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return
}

// BeforeUpdate hook: Task güncellenmeden önce çalışır.
func (t *Task) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now()
	return
}
