package db

import (
	"Vnollx/utils"
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type Notification struct {
	gorm.Model
	Title      string `gorm:"not null" json:"title,omitempty"`
	Ascription int64  `gorm:"not null" json:"ascription,omitempty"` //属于谁
	Content    string `gorm:"not null" json:"content,omitempty"`
	CreateTime string `gorm:"not null" json:"create_time,omitempty"`
}

func GetNotificationsByUserId(id int64) ([]*Notification, error) {
	var result []*Notification
	if err := DB.Where("ascription = ?", id).Find(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
func GetNotificationById(id int64) (*Notification, error) {
	var notication *Notification
	if err := DB.Where("ID = ?", id).First(&notication).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return notication, nil
}
func DeleteNotification(notification *Notification) error {
	err := DB.Unscoped().Delete(notification).Error
	return err
}
func CreateNotification(title string, content string, ascription int64) error {
	nowtime, err := utils.GetNowTime()
	if err != nil {
		return err
	}
	notification := &Notification{
		Title:      title,
		Ascription: ascription,
		Content:    content,
		CreateTime: nowtime,
	}
	if err = DB.Create(notification).Error; err != nil {
		return err
	}
	return nil
}
func SendNotificationIfSpaceNotEnough() {
	for {
		time.Sleep(time.Hour)
		var user []*User
		err := DB.Where("space< ?", 1073741824).Find(&user).Error
		if err != nil {
			log.Println("获取用户信息失败")
		}
		for _, u := range user {
			err = CreateNotification("系统通知", "云盘空间已不足1GB，请注意", int64(u.ID))
			if err != nil {
				log.Println("发送系统通知失败")
			}
		}
	}
}
