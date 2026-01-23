package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`

	AuthorID uint `gorm:"not null"`
	PostID   uint `gorm:"not null"`

	Author *User `gorm:"foreignKey:AuthorID"`
	Post   *Post `gorm:"foreignKey:PostID"`
}
