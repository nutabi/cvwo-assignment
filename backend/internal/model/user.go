package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	PHC      string `gorm:"not null"`

	AvatarURL *string
	Bio       *string

	Posts    []*Post    `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	Comments []*Comment `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
}
