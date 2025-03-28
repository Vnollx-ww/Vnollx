package main

import (
	"Vnollx/dal/db"
	"Vnollx/kitex_gen/base"
	user "Vnollx/kitex_gen/user"
	"Vnollx/pkg/middlerware"
	"Vnollx/utils"
	"context"
	"fmt"
	"go.uber.org/zap"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// UserLogin implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserLogin(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	resp = new(user.UserLoginResponse)
	usr, err := db.GetUserByPhone(req.Phone)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if usr == nil {
		resp.Base = utils.WrapWithLogicError("用户不存在")
		return resp, nil
	}
	isValid := utils.VerifyPassword(usr.Password, usr.Salt, req.Password)
	if !isValid {
		resp.Base = utils.WrapWithLogicError("密码错误")
		return resp, nil
	}
	token, err := middlerware.NewToken(int64(usr.ID))
	if err != nil {
		logger.Error("Token生成失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Token = token
	resp.Base = utils.WrapWithSuccess("登录成功")
	return resp, nil
}

// UserResgiter implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserResgiter(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	resp = new(user.UserRegisterResponse)
	usr, err := db.GetUserByPhone(req.Phone)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if usr != nil {
		resp.Base = utils.WrapWithLogicError("手机号已存在")
		return resp, nil
	}
	salt, err := utils.GenerateSalt(16)
	if err != nil {
		logger.Error("加密失败：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	hashedPassword := utils.HashPassword(req.Password, salt)
	err = db.CreateUser(req.Name, req.Phone, hashedPassword, salt, req.Avatar)
	if err != nil {
		logger.Error("注册失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("注册成功")
	return resp, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	resp = new(user.GetUserInfoResponse)
	usr, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if usr == nil {
		logger.Error("用户不存在", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("手机号已存在")
		return resp, nil
	}
	resp.Userinfo = new(base.UserInfo)
	resp.Userinfo.Phone = usr.Phone
	resp.Userinfo.Name = usr.Name
	resp.Userinfo.Space = fmt.Sprintf("%d", usr.Space)
	resp.Userinfo.Avatar = usr.Avatar
	resp.Userinfo.Id = int64(usr.ID)
	resp.Base = utils.WrapWithSuccess("获取信息成功")
	return resp, nil
}
func (s *UserServiceImpl) UpdateUserInfo(ctx context.Context, req *user.UpdateUserInfoRequest) (resp *user.UpdateUserInfoResponse, err error) {
	resp = new(user.UpdateUserInfoResponse)
	usr, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if usr == nil {
		logger.Error("用户不存在", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("手机号已存在")
		return resp, nil
	}
	err = db.UpdateUserInfo(usr, req.Name, req.Password, req.Phone, req.Avatar)
	if err != nil {
		logger.Error("更新用户信息失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("修改信息成功")
	return resp, nil
}
