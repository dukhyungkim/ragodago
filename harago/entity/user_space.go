package entity

import "time"

type UserSpace struct {
	ID        uint
	Name      string `gorm:"size:16;not null"`
	Email     string `gorm:"size:64;not null;unique"`
	Space     string `gorm:"size:32;not null;unique"`
	CreatedAt time.Time
}
