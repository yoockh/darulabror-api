package models

import "gorm.io/datatypes"

type Article struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Content   datatypes.JSON `gorm:"type:jsonb;not null" json:"content"`
	Author    string         `gorm:"not null" json:"author"`
	Status    string         `gorm:"not null;default:'draft'" json:"status"` // draft / published
	CreatedAt int64          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64          `gorm:"autoUpdateTime" json:"updated_at"`
}
