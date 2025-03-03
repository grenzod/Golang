package models

type Dialog struct {
	ID      int64  `gorm:"primaryKey"`
	Lang    string `gorm:"type:varchar(2);not null"`
	Content string `gorm:"type:text;not null"`
	Words   []Word `gorm:"many2many:word_dialog;"`
}

func (Dialog) TableName() string {
	return "dialog"
}

type Word struct {
	ID        int64    `gorm:"primaryKey"`
	Lang      string   `gorm:"type:varchar(2);not null"`
	Content   string   `gorm:"type:text;not null"`
	Translate string   `gorm:"type:text;not null"`
	Dialogs   []Dialog `gorm:"many2many:word_dialog;"`
}

func (Word) TableName() string {
	return "word"
}
