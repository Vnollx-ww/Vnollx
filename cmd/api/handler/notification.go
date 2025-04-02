package handler

import (
	"Vnollx/cmd/api/rpc"
	"Vnollx/kitex_gen/notification"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func GetNotificationsByUserId(ctx *gin.Context) {
	req := &notification.GetNotificationsByUserIdRequest{
		UserId: ctx.Value("uid").(int64),
	}
	resp, _ := rpc.GetNotificationsByUserId(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func DeleteNotification(ctx *gin.Context) {
	var reqbody struct {
		NotificationId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	nid, _ := strconv.ParseInt(reqbody.NotificationId, 10, 64)
	req := &notification.DeleteNotificationRequest{
		NotificationId: nid,
		UserId:         ctx.Value("uid").(int64),
	}
	resp, _ := rpc.DeleteNotification(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func SendNotification(ctx *gin.Context) {
	var reqbody struct {
		ToUserId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	touserid, _ := strconv.ParseInt(reqbody.ToUserId, 10, 64)
	req := &notification.SendNotificationRequest{
		TouserId: touserid,
		UserId:   ctx.Value("uid").(int64),
	}
	resp, _ := rpc.SendNotification(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
