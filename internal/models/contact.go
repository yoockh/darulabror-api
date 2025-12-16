package models

type Contact struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string `gorm:"not null" json:"email"`
	Subject   string `gorm:"not null" json:"subject"`
	Message   string `gorm:"type:text;not null" json:"message"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
}
