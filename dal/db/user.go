package db

import (
	"Vnollx/utils"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	Name        string `gorm:"varchar(40);not null" json:"name,omitempty"`
	Phone       string `gorm:"unique;varchar(20);not null" json:"phone,omitempty"`
	Password    []byte `gorm:"varchar(256);not null" json:"password,omitempty"`
	Avatar      string `gorm:"varchar(256);not null" json:"avatar,omitempty"`
	Salt        []byte `gorm:"varchar(256);not null" json:"-"`
	Space       int64  `gorm:"not null;default:10737418240" json:"space,omitempty"`
	Friend      string `gorm:"type:json" json:"friend,omitempty"`
	SendMessage string `gorm:"type:json" json:"send_message,omitempty"` // 设置默认值为空 JSON 对象
}

func CreateUser(name string, phone string, password []byte, salt []byte, postfix string) (string, error) {
	usr := &User{Name: name, Phone: phone, Password: password, Salt: salt, Friend: "{}", SendMessage: "{}"}
	err := DB.Create(usr).Error
	if err != nil {
		return "", err
	}
	usr.Avatar = fmt.Sprintf("http://106.54.223.38:9000/user/%d.%s", usr.ID, postfix)
	log.Println(usr.Avatar)
	err = DB.Save(usr).Error
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", usr.ID), err
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
func UpdateUserInfo(usr *User, name string, phone string) error {
	usr.Name = name
	if phone != usr.Phone {
		usr.Phone = phone
		err := CreateNotification("系统通知", "个人信息修改成功", int64(usr.ID))
		if err != nil {
			return err
		}
	}
	if err := DB.Save(usr).Error; err != nil {
		return err
	}
	return nil
}
func UpdateUserPassword(usr *User, password string) error {
	if password != "" {
		hashedPassword := utils.HashPassword(password, usr.Salt)
		usr.Password = hashedPassword
	}
	if err := DB.Save(usr).Error; err != nil {
		return err
	}
	return nil
}
func AddFriend(usra *User, usrb *User) error {
	var usrafriend map[int64]int
	var usrbfriend map[int64]int
	err := json.Unmarshal([]byte(usra.Friend), &usrafriend)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(usrb.Friend), &usrbfriend)
	if err != nil {
		return err
	}
	_, exists := usrafriend[int64(usrb.ID)]
	if exists {
		return fmt.Errorf("已经是好友了，无需再次添加")
	}
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 添加好友关系
	usrafriend[int64(usrb.ID)] = 1 // 将 usrb 添加到 usra 的好友列表
	usrbfriend[int64(usra.ID)] = 1 // 将 usra 添加到 usrb 的好友列表
	usraJSON, err := json.Marshal(usrafriend)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("序列化 用户A 的 Friend 错误: %v", err)
	}
	usrbJSON, err := json.Marshal(usrbfriend)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("序列化 用户b 的 Friend 错误: %v", err)
	}
	usra.Friend = string(usraJSON)
	usrb.Friend = string(usrbJSON)
	if err = tx.Save(&usra).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Save(&usrb).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func DeleteFriend(usra *User, usrb *User) error {
	var usrafriend map[int64]int
	var usrbfriend map[int64]int
	err := json.Unmarshal([]byte(usra.Friend), &usrafriend)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(usrb.Friend), &usrbfriend)
	if err != nil {
		return err
	}
	_, exists := usrafriend[int64(usrb.ID)]
	if !exists {
		return fmt.Errorf("已经不再是好友了，无需再次删除")
	}
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	delete(usrafriend, int64(usrb.ID)) // 从 usra 的好友列表中删除 usrb
	delete(usrbfriend, int64(usra.ID)) // 从 usrb 的好友列表中删除 usra
	usraJSON, err := json.Marshal(usrafriend)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("序列化 用户A 的 Friend 错误: %v", err)
	}
	usrbJSON, err := json.Marshal(usrbfriend)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("序列化 用户B 的 Friend 错误: %v", err)
	}
	usra.Friend = string(usraJSON)
	usrb.Friend = string(usrbJSON)

	// 删除 SendMessage 中与对方的聊天记录
	var sendMessageA map[int64][][2]string
	var sendMessageB map[int64][][2]string
	err = json.Unmarshal([]byte(usra.SendMessage), &sendMessageA)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("反序列化 用户A 的 SendMessage 错误: %v", err)
	}
	err = json.Unmarshal([]byte(usrb.SendMessage), &sendMessageB)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("反序列化 用户B 的 SendMessage 错误: %v", err)
	}

	// 删除与对方的聊天记录
	delete(sendMessageA, int64(usrb.ID)) // 删除与 usrb 的聊天记录
	delete(sendMessageB, int64(usra.ID)) // 删除与 usra 的聊天记录

	// 更新 SendMessage 字段
	usraSendMessageJSON, err := json.Marshal(sendMessageA)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("序列化 用户A 的 SendMessage 错误: %v", err)
	}
	usrbSendMessageJSON, err := json.Marshal(sendMessageB)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("序列化 用户B 的 SendMessage 错误: %v", err)
	}
	usra.SendMessage = string(usraSendMessageJSON)
	usrb.SendMessage = string(usrbSendMessageJSON)

	// 保存更新后的数据
	if err = tx.Save(&usra).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Save(&usrb).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 提交事务
	return tx.Commit().Error
}
func DeleteMessage(usra *User, usrb *User) error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 删除 SendMessage 中与对方的聊天记录
	var sendMessageA map[int64][][2]string
	err := json.Unmarshal([]byte(usra.SendMessage), &sendMessageA)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("反序列化 用户A 的 SendMessage 错误: %v", err)
	}
	// 删除与对方的聊天记录
	delete(sendMessageA, int64(usrb.ID)) // 删除与 usrb 的聊天记录

	// 更新 SendMessage 字段
	usraSendMessageJSON, err := json.Marshal(sendMessageA)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("序列化 用户A 的 SendMessage 错误: %v", err)
	}
	usra.SendMessage = string(usraSendMessageJSON)
	if err = tx.Save(&usra).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func SendMessage(usr *User, tousr *User, msg string) error {
	var Sendmessage map[int64][][2]string
	err := json.Unmarshal([]byte(usr.SendMessage), &Sendmessage)
	if err != nil {
		return err
	}
	currentTime, err := utils.GetNowTime()
	if err != nil {
		return err
	}
	messageItem := [2]string{msg, currentTime}
	//发送方添加发送消息到数据库
	if _, exists := Sendmessage[int64(tousr.ID)]; !exists {
		Sendmessage[int64(tousr.ID)] = [][2]string{}
	}
	Sendmessage[int64(tousr.ID)] = append(Sendmessage[int64(tousr.ID)], messageItem)
	sendmessage, err := json.Marshal(Sendmessage)
	if err != nil {
		return err
	}
	usr.SendMessage = string(sendmessage)
	err = DB.Save(&usr).Error
	if err != nil {
		return err
	}
	return nil
}
func IsFriendJudge(usr *User, ID int64) (bool, error) {
	var usrfriend map[int64]int
	err := json.Unmarshal([]byte(usr.Friend), &usrfriend)
	if err != nil {
		return false, err
	}
	_, ok := usrfriend[ID]
	if ok {
		return true, nil
	} else {
		return false, nil
	}
}
func GetUsersByName(name string) ([]*User, error) {
	var users []*User
	// 使用MySQL的CAST函数将ID转换为字符串
	if err := DB.Where("name LIKE ? OR CAST(id AS CHAR) = ?", "%"+name+"%", name).Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return users, nil
}
