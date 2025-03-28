package db

import "log"

type FolderRelation struct {
	Pid uint64 `gorm:"primaryKey;not null"` // Pid 为 bigint unsigned 类型，作为复合主键的一部分
	Sid uint64 `gorm:"primaryKey;not null"` // Sid 为 bigint unsigned 类型，作为复合主键的一部分

	// 外键约束
	FolderPid Folder `gorm:"foreignKey:Pid;references:ID;onDelete:CASCADE"`
	FolderSid Folder `gorm:"foreignKey:Sid;references:ID;onDelete:CASCADE"`
}

func CreateFolderRelation(pid uint64, sid uint64) error {
	// 查找 Folder 表中的 Folder 记录
	var folderPid Folder
	var folderSid Folder
	if err := DB.First(&folderPid, pid).Error; err != nil {
		log.Println("FolderPid not found:", err)
		return err
	}
	if err := DB.First(&folderSid, sid).Error; err != nil {
		log.Println("FolderSid not found:", err)
		return err
	}
	folderRelation := FolderRelation{
		Pid:       pid,
		Sid:       sid,
		FolderPid: folderPid,
		FolderSid: folderSid,
	}
	if err := DB.Create(&folderRelation).Error; err != nil {
		log.Println("Error creating FolderRelation:", err)
		return err
	}

	return nil
}
