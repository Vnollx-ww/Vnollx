package handler

import (
	"Vnollx/cmd/api/rpc"
	"Vnollx/kitex_gen/file"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func UploadFile(ctx *gin.Context) {
	var reqbody struct {
		ParentFolderId string
		Name           string
		Size           string
		Source         string
		Postfix        string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	pfid, _ := strconv.ParseInt(reqbody.ParentFolderId, 10, 64)
	req := &file.UploadFileRequest{
		UserId:         ctx.Value("uid").(int64),
		ParentFolderId: pfid,
		Name:           reqbody.Name,
		Size:           reqbody.Size,
		Source:         reqbody.Source,
		Postfix:        reqbody.Postfix,
	}
	resp, _ := rpc.UploadFile(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func DeleteFile(ctx *gin.Context) {
	var reqbody struct {
		FileId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	fid, _ := strconv.ParseInt(reqbody.FileId, 10, 64)
	req := &file.DeleteFileRequest{
		FileId: fid,
		UserId: ctx.Value("uid").(int64),
	}
	resp, _ := rpc.DeleteFile(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func UpdateFileInfo(ctx *gin.Context) {
	var reqbody struct {
		FileId         string
		Name           string
		ParentFolderId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	fid, _ := strconv.ParseInt(reqbody.FileId, 10, 64)
	pfid, _ := strconv.ParseInt(reqbody.ParentFolderId, 10, 64)
	req := &file.UpdateFileInfoRequest{
		FileId:         fid,
		UserId:         ctx.Value("uid").(int64),
		Name:           reqbody.Name,
		ParentFolderId: pfid,
	}
	resp, _ := rpc.UpdateFileInfo(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}

func GetFileInfo(ctx *gin.Context) {
	var reqbody struct {
		FileId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	fid, _ := strconv.ParseInt(reqbody.FileId, 10, 64)
	req := &file.GetFileInfoRequest{
		FileId: fid,
	}
	resp, _ := rpc.GetFileInfo(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetFilesInfo(ctx *gin.Context) {
	var reqbody struct {
		ParentFolderId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	pfid, _ := strconv.ParseInt(reqbody.ParentFolderId, 10, 64)
	req := &file.GetFilesInfoRequest{
		ParentFolderId: pfid,
		UserId:         ctx.Value("uid").(int64),
	}
	resp, _ := rpc.GetFilesInfo(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetAllFiles(ctx *gin.Context) {
	var reqbody struct {
		Name string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &file.GetALLFilesRequest{
		Name:   reqbody.Name,
		UserId: ctx.Value("uid").(int64),
	}
	resp, _ := rpc.GetAllFiles(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetVideoFiles(ctx *gin.Context) {
	var reqbody struct {
		Name string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &file.GetVideoFilesRequest{
		Name:   reqbody.Name,
		UserId: ctx.Value("uid").(int64),
	}
	resp, _ := rpc.GetVideoFiles(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetMusicFiles(ctx *gin.Context) {
	var reqbody struct {
		Name string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &file.GetMusicFilesRequest{
		Name:   reqbody.Name,
		UserId: ctx.Value("uid").(int64),
	}
	resp, _ := rpc.GetMusicFiles(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetPictureFiles(ctx *gin.Context) {
	var reqbody struct {
		Name string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &file.GetPictureFilesRequest{
		Name:   reqbody.Name,
		UserId: ctx.Value("uid").(int64),
	}
	resp, _ := rpc.GetPictureFiles(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetDocumentFiles(ctx *gin.Context) {
	var reqbody struct {
		Name string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &file.GetDocumentFilesRequest{
		Name:   reqbody.Name,
		UserId: ctx.Value("uid").(int64),
	}
	resp, _ := rpc.GetDocumentFiles(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetOtherFiles(ctx *gin.Context) {
	var reqbody struct {
		Name string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &file.GetOtherFilesRequest{
		Name:   reqbody.Name,
		UserId: ctx.Value("uid").(int64),
	}
	resp, _ := rpc.GetOtherFiles(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetFilesByName(ctx *gin.Context) {
	var reqbody struct {
		Name string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &file.GetFilesByNameRequest{
		Name:   reqbody.Name,
		UserId: ctx.Value("uid").(int64),
	}
	resp, _ := rpc.GetFilesByName(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func CreateShare(ctx *gin.Context) {
	var reqbody struct {
		FileId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	fid, _ := strconv.ParseInt(reqbody.FileId, 10, 64)
	req := &file.CreateShareRequest{
		UserId: ctx.Value("uid").(int64),
		FileId: fid,
	}
	resp, _ := rpc.CreateShare(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetFileByCode(ctx *gin.Context) {
	var reqbody struct {
		Code string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &file.GetFileByCodeRequest{
		Code: reqbody.Code,
	}
	resp, _ := rpc.GetFileByCode(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
