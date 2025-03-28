package handler

import (
	"Vnollx/cmd/api/rpc"
	"Vnollx/dal/redis"
	"Vnollx/kitex_gen/user"
	captcha2 "Vnollx/pkg/captcha"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Login(ctx *gin.Context) {
	var reqbody struct {
		Phone    string
		Password string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &user.UserLoginRequest{
		Phone:    reqbody.Phone,
		Password: reqbody.Password,
	}
	resp, _ := rpc.Login(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func Register(ctx *gin.Context) {
	var reqbody struct {
		Name     string
		Phone    string
		Password string
		Avatar   string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &user.UserRegisterRequest{
		Name:     reqbody.Name,
		Phone:    reqbody.Phone,
		Password: reqbody.Password,
		Avatar:   reqbody.Avatar,
	}
	resp, _ := rpc.Register(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetUserinfo(ctx *gin.Context) {
	req := &user.GetUserInfoRequest{
		UserId: ctx.Value("uid").(int64),
	}
	resp, _ := rpc.GetUserInfo(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func UpdateUserInfo(ctx *gin.Context) {
	var reqbody struct {
		Name     string
		Phone    string
		Password string
		Avatar   string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &user.UpdateUserInfoRequest{
		UserId:   ctx.Value("uid").(int64),
		Name:     reqbody.Name,
		Phone:    reqbody.Phone,
		Avatar:   reqbody.Avatar,
		Password: reqbody.Password,
	}
	resp, _ := rpc.UpdateUserInfo(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetCaptcha(ctx *gin.Context) {
	captcha := captcha2.GenerateCaptcha(ctx)
	key := "captcha:" + ctx.ClientIP()
	err := redis.SetKey(context.Background(), key, captcha, 2*time.Minute)
	if err != nil {
		logger.Error("验证码存入redis失败", zap.Error(err))
	}
}
