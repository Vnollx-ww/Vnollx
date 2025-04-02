package handler

import (
	"Vnollx/cmd/api/rpc"
	"Vnollx/kitex_gen/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
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
		Postfix  string
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
		Postfix:  reqbody.Postfix,
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
		Name  string
		Phone string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &user.UpdateUserInfoRequest{
		UserId: ctx.Value("uid").(int64),
		Name:   reqbody.Name,
		Phone:  reqbody.Phone,
	}
	resp, _ := rpc.UpdateUserInfo(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func UpdatePassword(ctx *gin.Context) {
	var reqbody struct {
		OldPassword string
		NewPassword string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &user.UpdatePasswordRequest{
		OldPassword:  reqbody.OldPassword,
		NewPassword_: reqbody.NewPassword,
		UserId:       ctx.Value("uid").(int64),
	}
	resp, _ := rpc.UpdatePassword(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func UserLoginByCode(ctx *gin.Context) {
	var reqbody struct {
		Captcha string
		Phone   string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &user.UserLoginByCodeRequest{
		Captcha: reqbody.Captcha,
		Phone:   reqbody.Phone,
	}
	resp, _ := rpc.LoginByCode(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GenerateCaptcha(ctx *gin.Context) {
	var reqbody struct {
		Phone string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &user.GenerateCaptchaRequest{
		Phone: reqbody.Phone,
	}
	resp, _ := rpc.GenerateCaptcha(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func AddFriend(ctx *gin.Context) {
	var reqbody struct {
		ToUserId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	touserid, _ := strconv.ParseInt(reqbody.ToUserId, 10, 64)
	req := &user.AddFriendRequest{
		UserId:   ctx.Value("uid").(int64),
		TouserId: touserid,
	}
	resp, _ := rpc.AddFriend(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func DeleteFriend(ctx *gin.Context) {
	var reqbody struct {
		ToUserId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	touserid, _ := strconv.ParseInt(reqbody.ToUserId, 10, 64)
	req := &user.DeleteFriendRequest{
		UserId:   ctx.Value("uid").(int64),
		TouserId: touserid,
	}
	resp, _ := rpc.DeleteFriend(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func SendMessage(ctx *gin.Context) {
	var reqbody struct {
		ToUserId string
		Content  string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	touserid, _ := strconv.ParseInt(reqbody.ToUserId, 10, 64)
	req := &user.SendMessageRequest{
		UserId:   ctx.Value("uid").(int64),
		TouserId: touserid,
		Content:  reqbody.Content,
	}
	resp, _ := rpc.SendMessage(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetFriendList(ctx *gin.Context) {
	req := &user.GetFriendListRequest{
		UserId: ctx.Value("uid").(int64),
	}
	resp, _ := rpc.GetFriendList(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetMessageList(ctx *gin.Context) {
	var reqbody struct {
		ToUserId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	touserid, _ := strconv.ParseInt(reqbody.ToUserId, 10, 64)
	req := &user.GetMessageListRequest{
		UserId:   ctx.Value("uid").(int64),
		TouserId: touserid,
	}
	resp, _ := rpc.GetMessageList(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func DeleteMessage(ctx *gin.Context) {
	var reqbody struct {
		ToUserId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	touserid, _ := strconv.ParseInt(reqbody.ToUserId, 10, 64)
	req := &user.DeleteMessageRequest{
		UserId:   ctx.Value("uid").(int64),
		TouserId: touserid,
	}
	resp, _ := rpc.DeleteMessage(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func GetUsersByName(ctx *gin.Context) {
	var reqbody struct {
		Name string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	req := &user.GetUsersByNameRequest{
		Name: reqbody.Name,
	}
	resp, _ := rpc.GetUsersByName(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
func IsFriend(ctx *gin.Context) {
	var reqbody struct {
		ToUserId string
	}
	if err := ctx.ShouldBind(&reqbody); err != nil {
		logger.Error("前后端数据绑定错误", zap.Error(err))
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	touserid, _ := strconv.ParseInt(reqbody.ToUserId, 10, 64)
	req := &user.IsFriendRequest{
		UserId:   ctx.Value("uid").(int64),
		TouserId: touserid,
	}
	resp, _ := rpc.IsFriend(ctx, req)
	ctx.JSON(http.StatusOK, resp)
}
