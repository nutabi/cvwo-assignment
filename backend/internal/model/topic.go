package model

import "gorm.io/gorm"

type Topic struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"`
	Description *string

	AuthorID uint `gorm:"not null"`

	Author *User   `gorm:"foreignKey:AuthorID"`
	Posts  []*Post `gorm:"constraint:OnDelete:CASCADE"`
}
