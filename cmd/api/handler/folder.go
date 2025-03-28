package handler

import (
	"Vnollx/cmd/api/rpc"
	"Vnollx/kitex_gen/folder"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

func CreateFolder(ctx *gin.Context) {
	var reqbody struct {
		Name           string
		ParentFolderId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	pfid, _ := strconv.ParseInt(reqbody.ParentFolderId, 10, 64)
	req := &folder.CreateFolderRequest{
		Name:           reqbody.Name,
		ParentFolderId: pfid,
		UserId:         ctx.Value("uid").(int64),
	}
	resp, _ := rpc.CreateFolder(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func DeleteFolder(ctx *gin.Context) {
	var reqbody struct {
		FolderId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Println(reqbody.FolderId)
	fid, _ := strconv.ParseInt(reqbody.FolderId, 10, 64)
	req := &folder.DeleteFolderRequest{
		FolderId: fid,
		UserId:   ctx.Value("uid").(int64),
	}
	resp, _ := rpc.DeleteFolder(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func UpdateFolderInfo(ctx *gin.Context) {
	var reqbody struct {
		Name           string
		FolderId       string
		ParentFolderId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	fid, _ := strconv.ParseInt(reqbody.FolderId, 10, 64)
	pfid, _ := strconv.ParseInt(reqbody.ParentFolderId, 10, 64)
	req := &folder.UpdateFolderInfoRequest{
		Name:           reqbody.Name,
		FolderId:       fid,
		UserId:         ctx.Value("uid").(int64),
		ParentFolderId: pfid,
	}
	resp, _ := rpc.UpdateFolderInfo(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetFolderInfo(ctx *gin.Context) {
	var reqbody struct {
		FolderId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	fid, _ := strconv.ParseInt(reqbody.FolderId, 10, 64)
	req := &folder.GetFolderInfoRequest{
		FolderId: fid,
	}
	resp, _ := rpc.GetFolderInfo(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetFoldersInfo(ctx *gin.Context) {
	var reqbody struct {
		ParentFolderId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	pfid, _ := strconv.ParseInt(reqbody.ParentFolderId, 10, 64)
	req := &folder.GetFoldersInfoRequest{
		ParentFolderId: pfid,
		UserId:         ctx.Value("uid").(int64),
	}
	resp, _ := rpc.GetFoldersInfo(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
