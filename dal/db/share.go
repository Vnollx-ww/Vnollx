package db

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type Share struct {
	gorm.Model
	Code       string    `gorm:"not null" json:"code,omitempty"`
	FileId     int64     `gorm:"not null" json:"file_id,omitempty"`
	Ascription int64     `gorm:"not null" json:"ascription,omitempty"`
	ExpiredAt  time.Time `gorm:"not null" json:"expired_at,omitempty"`
}

func CreateShare(code string, fileId int64, ascription int64) error {
	share := &Share{
		Code:       code,
		FileId:     fileId,
		Ascription: ascription,
		ExpiredAt:  time.Now().Add(24 * time.Hour), // 设置过期时间为当前时间 + 24 小时
	}
	err := DB.Create(share).Error
	if err != nil {
		return err
	}
	return nil
}
func CleanExpiredShares() {
	for {
		time.Sleep(time.Hour)
		if err := DB.Where("expired_at < ?", time.Now()).Delete(&Share{}).Error; err != nil {
			log.Println("删除过期分享:", err)
		}
	}
}
func GetShareByCode(code string) (*Share, error) {
	var share *Share
	err := DB.Where("code=?", code).First(&share).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return share, nil
}
