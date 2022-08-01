package repository

import (
	"harago/entity"
	"log"
)

func (db *DB) SaveSpace(userSpace *entity.UserSpace) error {
	return db.client.Save(userSpace).Error
}

func (db *DB) FindSpaces() ([]*entity.UserSpace, error) {
	var spaces []*entity.UserSpace

	if err := db.client.Find(&spaces).Error; err != nil {
		return nil, err
	}

	return spaces, nil
}

func (db *DB) DeleteSpace(email string) {
	if err := db.client.Where(&entity.UserSpace{Email: email}).Delete(&entity.UserSpace{}).Error; err != nil {
		log.Println(err)
	}
}
