package tables

import (
	"time"

	"gorm.io/gorm"
)

type Provider struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Url       string    `gorm:"not null" json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Provider) TableName() string {
	return "tb_providers"
}

func (p *Provider) BeforeCreate(tx *gorm.DB) (err error) {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return
}

func (p *Provider) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return
}
