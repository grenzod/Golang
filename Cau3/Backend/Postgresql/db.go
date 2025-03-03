package db

import (
	"log"
	models "project/Model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InsertDialog(db *gorm.DB, dialog *models.Dialog) error {
	if err := db.Create(dialog).Error; err != nil {
		return err
	}
	return nil
}

func InsertWord(db *gorm.DB, word models.Word) error {
	if err := db.Create(&word).Error; err != nil {
		return err
	}
	return nil
}

func GetDialogWithWords(db *gorm.DB, id int64) {
	var dialog models.Dialog
	if err := db.Preload("Words").First(&dialog, id).Error; err != nil {
		log.Fatal("Error querying dialog:", err)
	}
	log.Printf("Dialog: %+v", dialog)
	log.Printf("Words: %+v", dialog.Words)
}

func AccessDb() *gorm.DB {
	dsn := "host=localhost user=postgres password=175003 dbname=Golang port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting:", err)
		return nil
	}
	return db
}
