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
