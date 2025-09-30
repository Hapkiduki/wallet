package domain

import "time"

type User struct {
	ID        string    `gorm:"type:uuid;primary_key"`
	Username  string    `gorm:"type:varchar(255);unique;not null"`
	Name      string    `gorm:"type:varchar(255);not null"`
	DNI       string    `gorm:"type:varchar(255);unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
