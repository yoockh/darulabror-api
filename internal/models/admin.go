package models

type Role string

var (
	Admins     Role = "admin"
	Superadmin Role = "superadmin"
)

type Admin struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string `gorm:"unique;not null" json:"username"`
	Password  string `gorm:"not null" json:"password"`
	Role      Role   `gorm:"not null;default:'admin'" json:"role"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
}
