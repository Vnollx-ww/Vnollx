package main

import (
	"Vnollx/dal/db"
	"Vnollx/kitex_gen/base"
	file "Vnollx/kitex_gen/file"
	"Vnollx/pkg/minio"
	"Vnollx/utils"
	"context"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

// FileServiceImpl implements the last service interface defined in the IDL.
type FileServiceImpl struct{}

// UploadFile implements the FileServiceImpl interface.
func (s *FileServiceImpl) UploadFile(ctx context.Context, req *file.UploadFileRequest) (resp *file.UploadFileResponse, err error) {
	resp = new(file.UploadFileResponse)
	size, _ := strconv.ParseInt(req.Size, 10, 64)
	usr, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	ok, err := db.GetFileByNameAndParentFolderId(req.Name, req.ParentFolderId)
	if err != nil {
		logger.Error("查询文件失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if ok == false {
		logger.Error("文件名称重复：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("文件名称重复")
		return resp, nil
	}
	ID, err := db.UploadFile(req.Name, req.Source, req.ParentFolderId, size, req.Postfix, req.UserId, usr)
	if err != nil {
		logger.Error("上传文件失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.FileId = fmt.Sprintf("%d", ID)
	resp.Base = utils.WrapWithSuccess("上传文件至云盘成功")
	return resp, nil
}

// DeleteFile implements the FileServiceImpl interface.
func (s *FileServiceImpl) DeleteFile(ctx context.Context, req *file.DeleteFileRequest) (resp *file.DeleteFileResponse, err error) {
	resp = new(file.DeleteFileResponse)
	f, err := db.GetFileById(req.FileId)
	if err != nil {
		logger.Error("查询文件失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if f == nil {
		logger.Error("文件不存在或已被删除：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("文件不存在或已被删除")
		return resp, nil
	}
	usr, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	err = db.DeleteFile(f, usr)
	if err != nil {
		logger.Error("删除文件失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	err = minio.DeleteFileMinio(ctx, int64(f.ID), f.Postfix)
	if err != nil {
		logger.Error("从云盘删除文件失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("从云盘删除文件成功")
	return resp, nil
}

// UpdateFileInfo implements the FileServiceImpl interface.
func (s *FileServiceImpl) UpdateFileInfo(ctx context.Context, req *file.UpdateFileInfoRequest) (resp *file.UpdateFileInfoResponse, err error) {
	resp = new(file.UpdateFileInfoResponse)
	f, err := db.GetFileById(req.FileId)
	if err != nil {
		logger.Error("查询文件失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if f == nil {
		logger.Error("文件不存在或已被删除：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("文件不存在或已被删除")
		return resp, nil
	}
	ok, err := db.GetFileByNameAndParentFolderId(req.Name, req.ParentFolderId)
	if err != nil {
		logger.Error("查询文件失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if ok == false {
		logger.Error("文件名称重复：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("文件名称重复")
		return resp, nil
	}
	err = db.UpdateFileInfo(f, req.Name)
	if err != nil {
		logger.Error("修改文件信息失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("更新文件属性成功")
	return resp, nil
}
func (s *FileServiceImpl) GetFileInfo(ctx context.Context, req *file.GetFileInfoRequest) (resp *file.GetFileInfoResponse, err error) {
	resp = new(file.GetFileInfoResponse)
	resp.File = new(base.FileInfo)
	f, err := db.GetFileById(req.FileId)
	if err != nil {
		logger.Error("查询文件失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if f == nil {
		logger.Error("文件不存在或已被删除：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("文件不存在或已被删除")
		return resp, nil
	}
	resp.File.Size = fmt.Sprintf("%d", f.Size)
	resp.File.Name = f.FileName
	resp.File.Postfix = f.Postfix
	resp.File.FileId = fmt.Sprintf("%d", f.ID)
	resp.File.ParentFolderId = fmt.Sprintf("%d", f.ParentFolderId)
	resp.File.UploadTime = f.UploadTime
	resp.Base = utils.WrapWithSuccess("获取文件信息成功")
	return resp, nil
}

// GetFilesInfo implements the FileServiceImpl interface.
func (s *FileServiceImpl) GetFilesInfo(ctx context.Context, req *file.GetFilesInfoRequest) (resp *file.GetFilesInfoResponse, err error) {
	resp = new(file.GetFilesInfoResponse)
	resp.Files = []*base.FileInfo{}
	files, err := db.GetFilesByUserIdAndParentFolderId(req.UserId, req.ParentFolderId)
	if err != nil {
		logger.Error("查询文件失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	for _, fi := range files {
		f := new(base.FileInfo)
		f.Name = fi.FileName
		f.Size = fmt.Sprintf("%d", fi.Size)
		f.Postfix = fi.Postfix
		f.UploadTime = fi.UploadTime
		f.FileId = fmt.Sprintf("%d", fi.ID)
		f.ParentFolderId = fmt.Sprintf("%d", fi.ParentFolderId)
		resp.Files = append(resp.Files, f)
	}
	resp.Base = utils.WrapWithSuccess("获取文件列表成功")
	return resp, nil
}
