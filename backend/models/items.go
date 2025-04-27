package models

import "gorm.io/gorm"

type Items struct {
	ID       uint    `gorm:"primary key;autoIncrement" json:"id"`
	Name     *string `json:"name"`
	Quantity *int    `json:"quantity"`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Items{})
	return err
}
