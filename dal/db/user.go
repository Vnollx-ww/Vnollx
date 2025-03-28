package db

import (
	"Vnollx/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"varchar(40);not null" json:"name,omitempty"`
	Phone    string `gorm:"unique;varchar(20);not null" json:"phone,omitempty"`
	Password []byte `gorm:"varchar(256);not null" json:"password,omitempty"`
	Avatar   string `gorm:"varchar(256);not null" json:"avatar,omitempty"`
	Salt     []byte `gorm:"varchar(256);not null" json:"-"`
	Space    int64  `gorm:"not null;default:10737418240" json:"space,omitempty"`
}

func CreateUser(name string, phone string, password []byte, salt []byte, avatar string) error {
	usr := &User{Name: name, Phone: phone, Password: password, Salt: salt, Avatar: avatar}
	err := DB.Create(usr).Error
	return err
}
func GetUserByPhone(phone string) (*User, error) {
	user := new(User)
	if err := DB.Where("phone = ?", phone).First(&user).Error; err == nil {
		return user, nil
	} else if user.ID == 0 {
		return nil, nil
	} else {
		return nil, err
	}
}
func GetUserById(id int64) (*User, error) {
	user := new(User)
	if err := DB.Where("ID = ?", id).First(&user).Error; err == nil {
		return user, nil
	} else if user.ID == 0 {
		return nil, nil
	} else {
		return nil, err
	}
}
func UpdateUserInfo(usr *User, name string, password string, phone string, avatar string) error {
	if name != "" {
		usr.Name = name
	}
	if password != "" {
		hashedPassword := utils.HashPassword(password, usr.Salt)
		usr.Password = hashedPassword
	}
	if phone != "" {
		usr.Phone = phone
	}
	if avatar != "" {
		usr.Avatar = avatar
	}
	if err := DB.Save(usr).Error; err != nil {
		return err
	}
	return nil
}
func UpdateUserSpace(usr *User, space int64) error {
	usr.Space = space
	if err := DB.Save(usr).Error; err != nil {
		return err
	}
	return nil
}
