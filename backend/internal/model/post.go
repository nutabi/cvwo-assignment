package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"type:text;not null"`

	AuthorID uint `gorm:"not null"`
	TopicID  uint `gorm:"not null"`

	Author   *User      `gorm:"foreignKey:AuthorID"`
	Topic    *Topic     `gorm:"foreignKey:TopicID"`
	Comments []*Comment `gorm:"constraint:OnDelete:CASCADE"`
}
