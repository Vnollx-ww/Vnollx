package main

import (
	"Vnollx/dal/db"
	"Vnollx/kitex_gen/base"
	folder "Vnollx/kitex_gen/folder"
	"Vnollx/utils"
	"context"
	"fmt"
	"go.uber.org/zap"
)

// FolderServiceImpl implements the last service interface defined in the IDL.
type FolderServiceImpl struct{}

// CreateFolder implements the FolderServiceImpl interface.
func (s *FolderServiceImpl) CreateFolder(ctx context.Context, req *folder.CreateFolderRequest) (resp *folder.CreateFolderResponse, err error) {
	resp = new(folder.CreateFolderResponse)
	ok, err := db.GetFolderByNameAndParentFolderId(req.Name, req.ParentFolderId)
	if err != nil {
		logger.Error("查询文件夹失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if ok == false {
		logger.Error("文件夹名称重复：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("文件夹名称重复")
		return resp, nil
	}
	err = db.CreateFolder(req.Name, req.ParentFolderId, req.UserId)
	if err != nil {
		logger.Error("创建文件夹失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("创建文件夹成功")
	return resp, nil
}

// DeleteFolder implements the FolderServiceImpl interface.
func (s *FolderServiceImpl) DeleteFolder(ctx context.Context, req *folder.DeleteFolderRequest) (resp *folder.DeleteFolderResponse, err error) {
	resp = new(folder.DeleteFolderResponse)
	f, err := db.GetFolderById(req.FolderId)
	if err != nil {
		logger.Error("查询文件夹失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if f == nil {
		logger.Error("文件夹不存在或已被删除：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("文件夹不存在或已被删除")
		return resp, nil
	}
	usr, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	_, err = db.DeleteFolder(f, usr)
	if err != nil {
		logger.Error("删除文件夹失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("删除文件夹成功")
	return resp, nil
}

// UpdateFolderInfo implements the FolderServiceImpl interface.
func (s *FolderServiceImpl) UpdateFolderInfo(ctx context.Context, req *folder.UpdateFolderInfoRequest) (resp *folder.UpdateFolderInfoResponse, err error) {
	resp = new(folder.UpdateFolderInfoResponse)
	f, err := db.GetFolderById(req.FolderId)
	if err != nil {
		logger.Error("查询文件夹失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if f == nil {
		logger.Error("文件夹不存在或已被删除：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("文件夹不存在或已被删除")
		return resp, nil
	}
	ok, err := db.GetFolderByNameAndParentFolderId(req.Name, req.ParentFolderId)
	if err != nil {
		logger.Error("查询文件夹失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if ok == false {
		logger.Error("文件夹名称重复：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("文件夹名称重复")
		return resp, nil
	}
	err = db.UpdateFolderInfo(f, req.Name)
	if err != nil {
		logger.Error("修改文件夹信息失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("更新文件夹属性成功")
	return resp, nil
}
func (s *FolderServiceImpl) GetFolderInfo(ctx context.Context, req *folder.GetFolderInfoRequest) (resp *folder.GetFolderInfoResponse, err error) {
	resp = new(folder.GetFolderInfoResponse)
	resp.Folder = new(base.FolderInfo)
	f, err := db.GetFolderById(req.FolderId)
	if err != nil {
		logger.Error("查询文件夹失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if f == nil {
		logger.Error("文件夹不存在或已被删除：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("文件夹不存在或已被删除")
		return resp, nil
	}
	resp.Folder.Name = f.FolderName
	resp.Folder.UploadTime = f.UploadTime
	resp.Folder.FolderId = fmt.Sprintf("%d", f.ID)
	resp.Folder.ParentFolderId = fmt.Sprintf("%d", f.ParentFolderId)
	resp.Base = utils.WrapWithSuccess("获取文件夹信息成功")
	return resp, nil
}

// GetFoldersInfo implements the FolderServiceImpl interface.
func (s *FolderServiceImpl) GetFoldersInfo(ctx context.Context, req *folder.GetFoldersInfoRequest) (resp *folder.GetFoldersInfoResponse, err error) {
	resp = new(folder.GetFoldersInfoResponse)
	resp.Folders = []*base.FolderInfo{}
	folders, err := db.GetFoldersByUserIdAndParentFolderId(req.UserId, req.ParentFolderId)
	if err != nil {
		logger.Error("查询文件夹失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	for _, fo := range folders {
		f := new(base.FolderInfo)
		f.Name = fo.FolderName
		f.UploadTime = fo.UploadTime
		f.ParentFolderId = fmt.Sprintf("%d", fo.ParentFolderId)
		f.FolderId = fmt.Sprintf("%d", fo.ID)
		resp.Folders = append(resp.Folders, f)
	}
	resp.Base = utils.WrapWithSuccess("获取文件夹列表成功")
	return resp, nil
}
