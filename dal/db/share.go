package db

import "gorm.io/gorm"

type Share struct {
	gorm.Model
	Code     string `gorm:"not null" json:"code,omitempty"`
	FileId   int64  `gorm:"not null" json:"file_id,omitempty"`
	Username string `gorm:"not null" json:"username,omitempty"`
	Hash     string `gorm:"not null" json:"hash,omitempty"`
}
