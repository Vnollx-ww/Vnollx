package main

import (
	"Vnollx/dal/db"
	"Vnollx/dal/redis"
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
	ok, err := db.GetFileByNameAndParentFolderId(req.Name, req.ParentFolderId, req.UserId)
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
	cacheKey := fmt.Sprintf("userinfo:%d", req.UserId)
	err = redis.DelKey(ctx, cacheKey) // 删除缓存
	if err != nil {
		logger.Error("删除缓存失败：", zap.Error(err))
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
	cacheKey := fmt.Sprintf("userinfo:%d", req.UserId)
	err = redis.DelKey(ctx, cacheKey) // 删除缓存
	if err != nil {
		logger.Error("删除缓存失败：", zap.Error(err))
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
	ok, err := db.GetFileByNameAndParentFolderId(req.Name, req.ParentFolderId, req.UserId)
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

// GetAllFiles implements the FileServiceImpl interface.
func (s *FileServiceImpl) GetAllFiles(ctx context.Context, req *file.GetALLFilesRequest) (resp *file.GetAllFilesResponse, err error) {
	resp = new(file.GetAllFilesResponse)
	resp.Files = []*base.FileInfo{}
	files, err := db.GetFilesByUserId(req.Name, req.UserId)
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

// GetVideoFiles implements the FileServiceImpl interface.
func (s *FileServiceImpl) GetVideoFiles(ctx context.Context, req *file.GetVideoFilesRequest) (resp *file.GetVideoFilesResponse, err error) {
	resp = new(file.GetVideoFilesResponse)
	resp.Files = []*base.FileInfo{}
	files, _, err := db.GetVideosByUserId(req.Name, req.UserId, 1)
	if err != nil {
		logger.Error("查询视频失败：", zap.Error(err))
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
	resp.Base = utils.WrapWithSuccess("获取视频列表成功")
	return resp, nil
}

// GetMusicFiles implements the FileServiceImpl interface.
func (s *FileServiceImpl) GetMusicFiles(ctx context.Context, req *file.GetMusicFilesRequest) (resp *file.GetMusicFilesResponse, err error) {
	resp = new(file.GetMusicFilesResponse)
	resp.Files = []*base.FileInfo{}
	files, _, err := db.GetMusicsByUserId(req.Name, req.UserId, 1)
	if err != nil {
		logger.Error("查询音频失败：", zap.Error(err))
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
	resp.Base = utils.WrapWithSuccess("获取音频列表成功")
	return resp, nil
}

// GetPictureFiles implements the FileServiceImpl interface.
func (s *FileServiceImpl) GetPictureFiles(ctx context.Context, req *file.GetPictureFilesRequest) (resp *file.GetPictureFilesResponse, err error) {
	resp = new(file.GetPictureFilesResponse)
	resp.Files = []*base.FileInfo{}
	files, _, err := db.GetPicturesByUserId(req.Name, req.UserId, 1)
	if err != nil {
		logger.Error("查询图片失败：", zap.Error(err))
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
	resp.Base = utils.WrapWithSuccess("获取图片列表成功")
	return resp, nil
}

// GetDocumentFiles implements the FileServiceImpl interface.
func (s *FileServiceImpl) GetDocumentFiles(ctx context.Context, req *file.GetDocumentFilesRequest) (resp *file.GetDocumentFilesResponse, err error) {
	resp = new(file.GetDocumentFilesResponse)
	resp.Files = []*base.FileInfo{}
	files, _, err := db.GetDocumentsByUserId(req.Name, req.UserId, 1)
	if err != nil {
		logger.Error("查询文档失败：", zap.Error(err))
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
	resp.Base = utils.WrapWithSuccess("获取文档列表成功")
	return resp, nil
}

// GetOtherFiles implements the FileServiceImpl interface.
func (s *FileServiceImpl) GetOtherFiles(ctx context.Context, req *file.GetOtherFilesRequest) (resp *file.GetOtherFilesResponse, err error) {
	resp = new(file.GetOtherFilesResponse)
	resp.Files = []*base.FileInfo{}
	files, _, err := db.GetOthersByUserId(req.Name, req.UserId, 1)
	if err != nil {
		logger.Error("查询其他文件失败：", zap.Error(err))
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
	resp.Base = utils.WrapWithSuccess("获取其他文件列表成功")
	return resp, nil
}

// GetFilesByName implements the FileServiceImpl interface.
func (s *FileServiceImpl) GetFilesByName(ctx context.Context, req *file.GetFilesByNameRequest) (resp *file.GetFilesByNameResponse, err error) {
	resp = new(file.GetFilesByNameResponse)
	resp.Files = []*base.FileInfo{}
	files, err := db.GetFilesByNameAndUserId(req.Name, req.UserId)
	if err != nil {
		logger.Error("搜索文件失败：", zap.Error(err))
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
	resp.Base = utils.WrapWithSuccess("获取搜索文件列表成功")
	return resp, nil
}

// CreateShare implements the FileServiceImpl interface.
func (s *FileServiceImpl) CreateShare(ctx context.Context, req *file.CreateShareRequest) (resp *file.CreateShareResponse, err error) {
	resp = new(file.CreateShareResponse)
	code := utils.GenerateUniqueString()
	err = db.CreateShare(code, req.FileId, req.UserId)
	if err != nil {
		logger.Error("分享文件失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("分享文件成功")
	resp.Code = code
	return resp, nil
}

// GetFileByCode implements the FileServiceImpl interface.
func (s *FileServiceImpl) GetFileByCode(ctx context.Context, req *file.GetFileByCodeRequest) (resp *file.GetFileByCodeResponse, err error) {
	resp = new(file.GetFileByCodeResponse)
	resp.File = new(base.FileInfo)
	share, err := db.GetShareByCode(req.Code)
	if err != nil {
		logger.Error("获取分享失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if share == nil {
		logger.Error("分享码已过期：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("分享码已过期")
		return resp, nil
	}
	f, err := db.GetFileById(share.FileId)
	if err != nil {
		logger.Error("查找文件失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
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
