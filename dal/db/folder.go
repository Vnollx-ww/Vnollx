package db

import (
	"Vnollx/utils"
	"errors"
	"gorm.io/gorm"
)

type Folder struct {
	gorm.Model
	FolderName     string `gorm:"not null" json:"file_folder_name,omitempty"` //文件夹名称
	ParentFolderId int64  `gorm:"not null" json:"parent_folder_id,omitempty"` //父文件夹id
	Ascription     int64  `gorm:"not null" json:"ascription,omitempty"`       //下载次数
	UploadTime     string `gorm:"not null" json:"upload_time,omitempty"`
}

func GetFolderCountByUserId(userId int64) (count int, err error) {
	var folders []*Folder
	if err = DB.Where("ascription=?", userId).Find(&folders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return len(folders), nil
}
func CreateFolder(name string, parent_folder_id int64, ascription int64) error {
	tx := DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	time, _ := utils.GetNowTime()
	folder := &Folder{
		FolderName:     name,
		ParentFolderId: parent_folder_id,
		Ascription:     ascription,
		UploadTime:     time,
	}
	err := tx.Create(folder).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
func DeleteFolder(folder *Folder, usr *User) (int64, error) {
	tx := DB.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var folders []*Folder
	if err := tx.Where("parent_folder_id = ?", folder.ID).Find(&folders).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	sum := int64(0)
	var files []*File
	if err := tx.Where("parent_folder_id = ?", folder.ID).Find(&files).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	if len(files) != 0 {
		for _, file := range files {
			sum += file.Size
		}
		if err := tx.Unscoped().Delete(&files).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	for _, f := range folders {
		num, err := DeleteFolder(f, nil)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		sum += num
	}
	if err := tx.Unscoped().Delete(&folder).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	if usr != nil {
		usr.Space += sum
		if err := DB.Save(usr).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}
	return sum, nil
}

func GetFolderById(folder_id int64) (*Folder, error) {
	folder := new(Folder)
	if err := DB.Where("ID = ?", folder_id).First(&folder).Error; err == nil {
		return folder, nil
	} else if folder.ID == 0 {
		return nil, nil
	} else {
		return nil, err
	}
}
func UpdateFolderInfo(folder *Folder, name string) error {
	folder.FolderName = name
	err := DB.Save(folder).Error
	return err
}
func GetFoldersByUserIdAndParentFolderId(userId int64, parent_folder_id int64) ([]*Folder, error) {
	var result []*Folder
	if err := DB.Where("ascription = ? AND parent_folder_id = ?", userId, parent_folder_id).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func GetFolderByNameAndParentFolderId(folderName string, parent_folder_id int64, user_id int64) (bool, error) {
	var folder *Folder
	if err := DB.Where("folder_name = ? And parent_folder_id = ? And ascription = ? ", folderName, parent_folder_id, user_id).First(&folder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
