package db

import (
	"Vnollx/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// 文件表
type File struct {
	gorm.Model
	FileName       string `gorm:"not null" json:"file_name,omitempty"`        //文件名
	FileHash       string `gorm:"not null" json:"file_hash,omitempty"`        //文件哈希值
	FilePath       string `gorm:"not null" json:"file_path,omitempty"`        //文件存储路径
	UploadTime     string `gorm:"not null" json:"upload_time,omitempty"`      //上传文件时间
	ParentFolderId int64  `gorm:"not null" json:"parent_folder_id,omitempty"` //父文件夹id
	Size           int64  `gorm:"not null" json:"size,omitempty"`             //文件大小
	Type           int    `gorm:"not null" json:"type,omitempty"`             //文件类型
	Postfix        string `gorm:"not null" json:"postfix,omitempty"`          //文件后缀
	Ascription     int64  `gorm:"not null" json:"ascription,omitempty"`       //属于谁
}

var extToNumber map[string]int = map[string]int{
	".txt":  1,
	".jpg":  2,
	".png":  3,
	".pdf":  4,
	".docx": 5,
	".mp4":  6,
	".mp3":  7,
	".zip":  8,
}

func GetVideosByUserId(name string, userId int64, option int) ([]*File, int, error) {
	var files []*File
	if err := DB.Where("file_name LIKE ? AND ascription= ? AND postfix= ?", "%"+name+"%", userId, "mp4").Find(&files).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	if option == 1 {
		return files, 0, nil
	} else {
		return nil, len(files), nil
	}
}
func GetPicturesByUserId(name string, userId int64, option int) ([]*File, int, error) {
	var files []*File
	if err := DB.Where("file_name LIKE ? AND ascription = ? AND postfix IN (?)", "%"+name+"%", userId, []string{"jpg", "png"}).Find(&files).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	if option == 1 {
		return files, 0, nil
	} else {
		return nil, len(files), nil
	}
}
func GetDocumentsByUserId(name string, userId int64, option int) ([]*File, int, error) {
	var files []*File
	if err := DB.Where("file_name LIKE ? AND ascription = ? AND postfix IN (?)", "%"+name+"%", userId, []string{"docx", "txt", "pdf"}).Find(&files).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	if option == 1 {
		return files, 0, nil
	} else {
		return nil, len(files), nil
	}
}
func GetMusicsByUserId(name string, userId int64, option int) ([]*File, int, error) {
	var files []*File
	if err := DB.Where("file_name LIKE ? AND ascription = ? AND postfix = ?", "%"+name+"%", userId, "mp3").Find(&files).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	if option == 1 {
		return files, 0, nil
	} else {
		return nil, len(files), nil
	}
}
func GetOthersByUserId(name string, userId int64, option int) ([]*File, int, error) {
	var files []*File
	if err := DB.Where("file_name LIKE ? AND ascription = ? AND postfix NOT IN ?", "%"+name+"%", userId, []string{"docx", "txt", "pdf", "mp4", "mp3", "jpg", "png"}).Find(&files).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	if option == 1 {
		return files, 0, nil
	} else {
		return nil, len(files), nil
	}
}
func UploadFile(name string, source string, parent_folder_id int64, size int64, postfix string, ascription int64, usr *User) (int64, error) {
	tx := DB.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取当前时间
	nowtime, err := utils.GetNowTime()
	if err != nil {
		return 0, err
	}
	number, _ := extToNumber[postfix]
	// 创建新的文件记录
	file := &File{
		FileName:       name,
		FileHash:       "hash", // 这里可以替换为计算出的文件hash
		FilePath:       source,
		UploadTime:     nowtime,
		ParentFolderId: parent_folder_id,
		Size:           size,
		Type:           number,
		Postfix:        postfix,
		Ascription:     ascription,
	}
	err = tx.Create(file).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	usr.Space -= size
	err = tx.Save(usr).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	file.FilePath = fmt.Sprintf("FileId=%d.%s", file.ID, postfix)
	err = tx.Save(file).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	err = tx.Commit().Error
	if err != nil {
		return 0, err
	}
	return int64(file.ID), nil
}

func GetFileById(id int64) (*File, error) {
	file := new(File)
	if err := DB.Where("ID = ?", id).First(&file).Error; err == nil {
		return file, nil
	} else if file.ID == 0 {
		return nil, nil
	} else {
		return nil, err
	}
}
func DeleteFile(file *File, usr *User) error {
	tx := DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	usr.Space += file.Size
	if err := tx.Save(usr).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Unscoped().Delete(file).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func UpdateFileInfo(file *File, name string) error {
	file.FileName = name
	err := DB.Save(file).Error
	return err
}
func GetFilesByUserId(name string, userId int64) ([]*File, error) {
	var files []*File
	if err := DB.Where("file_name LIKE ? AND ascription=?", "%"+name+"%", userId).Find(&files).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return files, nil
}
func GetFilesByUserIdAndParentFolderId(userId int64, parent_folder_id int64) ([]*File, error) {
	var result []*File
	if err := DB.Where("ascription = ? AND parent_folder_id = ?", userId, parent_folder_id).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func GetFileByNameAndParentFolderId(fileName string, parent_folder_id int64, user_id int64) (bool, error) {
	var file *File
	if err := DB.Where("file_name = ? And parent_folder_id = ? AND ascription=?", fileName, parent_folder_id, user_id).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
func GetFilesByNameAndUserId(fileName string, user_id int64) ([]*File, error) {
	var result []*File
	if err := DB.Where("file_name LIKE ? AND ascription = ?", "%"+fileName+"%", user_id).Find(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
