package repository

import (
	"fmt"
	"harago/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	client *gorm.DB
}

func NewPostgres(cfg *config.DB) (*DB, error) {
	const dnsTemplate = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s"
	dsn := fmt.Sprintf(dnsTemplate, cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port, cfg.Timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{client: db}, nil
}
