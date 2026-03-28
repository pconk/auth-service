package entity

import "time"

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"unique;not null;type:varchar(50)" json:"username"`
	PasswordHash string    `gorm:"not null;type:varchar(255)" json:"-"`
	Email        string    `gorm:"unique;not null;type:varchar(100)" json:"email"`
	Role         string    `gorm:"type:enum('admin','staff');default:'staff'" json:"role"`
	Position     string    `gorm:"type:varchar(100)" json:"position"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
