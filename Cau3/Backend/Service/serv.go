package service

import (
	"fmt"
	models "project/Model"
	db "project/Postgresql"

	"github.com/kataras/iris/v12"
)

func CreateData(words []models.Word, dialog models.Dialog) error {
	data := db.AccessDb()

	if err := db.InsertDialog(data, &dialog); err != nil {
		return fmt.Errorf("internal server error: %d", iris.StatusInternalServerError)
	}

	for _, w := range words {
		if err := db.InsertWord(data, w); err != nil {
			return fmt.Errorf("internal server error: %d", iris.StatusInternalServerError)
		}
	}

	data.Model(&dialog).Association("Words").Append(&words)

	return nil
}
