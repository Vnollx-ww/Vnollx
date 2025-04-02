package main

import (
	"Vnollx/dal/db"
	"Vnollx/dal/redis"
	"Vnollx/kitex_gen/base"
	user "Vnollx/kitex_gen/user"
	user2 "Vnollx/pkg/kafka/producer"
	"Vnollx/pkg/middlerware"
	"Vnollx/utils"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"
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
	id, err := db.CreateUser(req.Name, req.Phone, hashedPassword, salt, req.Postfix)
	if err != nil {
		logger.Error("注册失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("注册成功")
	resp.UserId = id
	return resp, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	resp = new(user.GetUserInfoResponse)
	resp.Userinfo = new(base.UserInfo)
	cachekey := fmt.Sprintf("userinfo:%d", req.UserId)
	u, err := redis.GetUserInfo(ctx, cachekey)
	if err == nil && u != nil {
		resp.Userinfo = u
		resp.Base = utils.WrapWithSuccess("获取信息成功")
		logger.Info("缓存信息获取成功")
		return resp, nil
	}
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
	resp.Userinfo.Phone = usr.Phone
	resp.Userinfo.Name = usr.Name
	resp.Userinfo.Space = fmt.Sprintf("%d", usr.Space)
	resp.Userinfo.Avatar = usr.Avatar
	resp.Userinfo.Id = int64(usr.ID)
	_, videocount, err := db.GetVideosByUserId("", (int64(usr.ID)), 0)
	if err != nil {
		logger.Error("查询视频失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	_, documentcount, err := db.GetDocumentsByUserId("", int64(usr.ID), 0)
	if err != nil {
		logger.Error("查询文档失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	_, picturecount, err := db.GetPicturesByUserId("", int64(usr.ID), 0)
	if err != nil {
		logger.Error("查询图片失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	_, musiccount, err := db.GetMusicsByUserId("", int64(usr.ID), 0)
	if err != nil {
		logger.Error("查询音频失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	_, othercount, err := db.GetOthersByUserId("", int64(usr.ID), 0)
	if err != nil {
		logger.Error("查询其他文件失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	foldercount, err := db.GetFolderCountByUserId(int64(usr.ID))
	if err != nil {
		logger.Error("查询文件夹失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Userinfo.Documentcount = fmt.Sprintf("%d", documentcount)
	resp.Userinfo.Musiccount = fmt.Sprintf("%d", musiccount)
	resp.Userinfo.Othercount = fmt.Sprintf("%d", othercount)
	resp.Userinfo.Picturecount = fmt.Sprintf("%d", picturecount)
	resp.Userinfo.Videocount = fmt.Sprintf("%d", videocount)
	resp.Userinfo.Foldercount = fmt.Sprintf("%d", foldercount)
	userdata, err := json.Marshal(resp.Userinfo)
	if err != nil {
		logger.Error("用户信息数据序列化失败：", zap.Error(err))
	}
	err = redis.SetKey(ctx, cachekey, userdata, 10*time.Minute)
	if err != nil {
		logger.Error("缓存设置失败：", zap.Error(err))
	}
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
	if usr.Phone != req.Phone {
		u, err := db.GetUserByPhone(req.Phone)
		if err != nil {
			logger.Error("查询用户失败", zap.Error(err))
			resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
			return resp, nil
		}
		if u != nil {
			logger.Error("手机号已存在", zap.Error(err))
			resp.Base = utils.WrapWithLogicError("手机号已存在")
			return resp, nil
		}
	}
	err = db.UpdateUserInfo(usr, req.Name, req.Phone)
	if err != nil {
		logger.Error("更新用户信息失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	cacheKey := fmt.Sprintf("userinfo:%d", req.UserId)
	err = redis.DelKey(ctx, cacheKey) // 删除缓存
	if err != nil {
		logger.Error("删除缓存失败：", zap.Error(err))
	}
	resp.Base = utils.WrapWithSuccess("修改信息成功")
	return resp, nil
}
func (s *UserServiceImpl) UpdatePassword(ctx context.Context, req *user.UpdatePasswordRequest) (resp *user.UpdatePasswordResponse, err error) {
	resp = new(user.UpdatePasswordResponse)
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
	isValid := utils.VerifyPassword(usr.Password, usr.Salt, req.OldPassword)
	if !isValid {
		resp.Base = utils.WrapWithLogicError("旧密码错误")
		return resp, nil
	}
	err = db.UpdateUserPassword(usr, req.NewPassword_)
	if err != nil {
		logger.Error("修改密码失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("修改密码成功")
	return resp, nil
}

// UserLoginByCode implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserLoginByCode(ctx context.Context, req *user.UserLoginByCodeRequest) (resp *user.UserLoginByCodeResponse, err error) {
	resp = new(user.UserLoginByCodeResponse)
	usr, err := db.GetUserByPhone(req.Phone)
	if err != nil {
		logger.Error("查询用户失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	key := "phone:" + req.Phone
	_, err = redis.IsCaptchaExists(ctx, key)
	if err != nil {
		logger.Error("验证码已过期", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("验证码已过期")
		return resp, nil
	}
	value, err := redis.GetKeyValue(ctx, key)
	if err != nil {
		logger.Error("获取redis验证码失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if value != req.Captcha {
		logger.Error("验证码错误，请重试", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("验证码错误，请重试")
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

// GenerateCaptcha implements the UserServiceImpl interface.
func (s *UserServiceImpl) GenerateCaptcha(ctx context.Context, req *user.GenerateCaptchaRequest) (resp *user.GenerateCaptchaResponse, err error) {
	resp = new(user.GenerateCaptchaResponse)
	captcha, err := utils.GenerateUniqueCaptcha(ctx, req.Phone)
	if err != nil {
		logger.Error("验证码生成失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	key := "phone:" + req.Phone
	err = redis.SetKey(context.Background(), key, captcha, 1*time.Minute)
	if err != nil {
		logger.Error("验证码存入redis失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("验证码已发送")
	return resp, nil
}

// AddFriend implements the UserServiceImpl interface.
func (s *UserServiceImpl) AddFriend(ctx context.Context, req *user.AddFriendRequest) (resp *user.AddFriendResponse, err error) {
	resp = new(user.AddFriendResponse)
	usra, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	usrb, err := db.GetUserById(req.TouserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	ok, err := db.IsFriendJudge(usra, req.TouserId)
	if err != nil {
		logger.Error("查询关系失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if ok == true {
		logger.Error("已是好友无需再添加：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("已是好友无需再添加")
		return resp, nil
	}
	err = db.AddFriend(usra, usrb)
	if err != nil {
		logger.Error("添加好友失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	cacheKey := "friendlist"
	err = redis.DelKey(ctx, cacheKey+fmt.Sprintf("%d", req.UserId)) // 删除缓存
	if err != nil {
		logger.Error("删除缓存失败：", zap.Error(err))
	}
	resp.Base = utils.WrapWithSuccess("添加好友成功")
	return resp, nil
}

// DeleteFriend implements the UserServiceImpl interface.
func (s *UserServiceImpl) DeleteFriend(ctx context.Context, req *user.DeleteFriendRequest) (resp *user.DeleteFriendResponse, err error) {
	resp = new(user.DeleteFriendResponse)
	usra, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	usrb, err := db.GetUserById(req.TouserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	ok, err := db.IsFriendJudge(usra, req.TouserId)
	if err != nil {
		logger.Error("查询关系失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if ok == false {
		logger.Error("已不是好友无需再删除：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("已不是好友无需再删除")
		return resp, nil
	}
	err = db.DeleteFriend(usra, usrb)
	if err != nil {
		logger.Error("删除好友失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	cacheKey := "friendlist"
	err = redis.DelKey(ctx, cacheKey+fmt.Sprintf("%d", req.UserId)) // 删除缓存
	if err != nil {
		logger.Error("删除缓存失败：", zap.Error(err))
	}
	resp.Base = utils.WrapWithSuccess("删除好友成功")
	return resp, nil
}

// SendMessage implements the UserServiceImpl interface.
func (s *UserServiceImpl) SendMessage(ctx context.Context, req *user.SendMessageRequest) (resp *user.SendMessageResponse, err error) {
	resp = new(user.SendMessageResponse)
	kafkaProducer, err1 := user2.NewKafkaProducer([]string{kafkaAddr}) //初始化kafka生产者
	var err2 error
	if err1 == nil {
		err2 = kafkaProducer.SendSendMessageEvent(req.UserId, req.Content, req.TouserId)
		if err2 == nil {
			resp.Base = utils.WrapWithSuccess("发送消息成功")
			return resp, nil
		}
	}
	if err1 != nil {
		logger.Error("kafka生产者创建失败：", zap.Error(err))
	}
	if err2 != nil {
		logger.Error("发送消息的请求创建成功，但消息发送到kafka失败：", zap.Error(err))
	}
	usra, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	usrb, err := db.GetUserById(req.TouserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	ok, err := db.IsFriendJudge(usra, req.TouserId)
	if err != nil {
		logger.Error("查询关系失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if ok == false {
		logger.Error("已不是好友不能发送消息：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("已不是好友不能发送消息")
		return resp, nil
	}
	err = db.SendMessage(usra, usrb, req.Content)
	if err != nil {
		logger.Error("发送消息失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	cacheKey := fmt.Sprintf("messagelist%d+%d", req.UserId, req.TouserId)
	err = redis.DelKey(ctx, cacheKey) // 删除缓存
	if err != nil {
		logger.Error("删除缓存失败：", zap.Error(err))
	}
	resp.Base = utils.WrapWithSuccess("发送消息成功")
	return resp, nil

}

// GetFriendList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFriendList(ctx context.Context, req *user.GetFriendListRequest) (resp *user.GetFriendListResponse, err error) {
	resp = new(user.GetFriendListResponse)
	resp.Friends = []*base.FriendInfo{}
	cacheKey := "friendlist"
	flist, err := redis.GetFriendList(ctx, cacheKey+fmt.Sprintf("%d", req.UserId))
	if err == nil && flist != nil {
		resp.Friends = flist
		resp.Base = utils.WrapWithSuccess("获取好友列表成功")
		return resp, nil
	}
	usr, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}

	var frlist map[int64]int
	err = json.Unmarshal([]byte(usr.Friend), &frlist)
	if err != nil {
		logger.Error("反序列化好友列表失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	for userid := range frlist {
		f := new(base.FriendInfo)
		u, err := db.GetUserById(userid)
		if err != nil {
			logger.Error("查询用户失败：", zap.Error(err))
			resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
			return resp, nil
		}
		f.Name = u.Name
		f.UserId = fmt.Sprintf("%d", u.ID)
		f.Avatar = u.Avatar
		resp.Friends = append(resp.Friends, f)
	}
	cacheddatabytes, err := json.Marshal(resp.Friends)
	if err != nil {
		logger.Error("好友列表数据序列化失败：", zap.Error(err))
	}
	err = redis.SetKey(ctx, cacheKey+fmt.Sprintf("%d", req.UserId), cacheddatabytes, 10*time.Minute)
	if err != nil {
		logger.Error("缓存设置失败：", zap.Error(err))
	}
	resp.Base = utils.WrapWithSuccess("获取好友列表成功")
	return resp, nil
}

// GetMessageList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetMessageList(ctx context.Context, req *user.GetMessageListRequest) (resp *user.GetMessageListResponse, err error) {
	resp = new(user.GetMessageListResponse)
	resp.Messages = []*base.Message{}
	cacheKey := fmt.Sprintf("messagelist%d+%d", req.UserId, req.TouserId)
	messagelist, err := redis.GetMessageListByFriend(ctx, cacheKey)
	if err == nil && messagelist != nil {
		resp.Messages = messagelist
		resp.Base = utils.WrapWithSuccess("获取好友聊天记录成功")
		return resp, nil
	}
	usr, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	tousr, err := db.GetUserById(req.TouserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	var smsg map[int64][][2]string
	smsg = make(map[int64][][2]string)
	err = json.Unmarshal([]byte(usr.SendMessage), &smsg)
	if err != nil {
		logger.Error("消息记录反序列化失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	var amsg map[int64][][2]string
	amsg = make(map[int64][][2]string)
	err = json.Unmarshal([]byte(tousr.SendMessage), &amsg)
	if err != nil {
		logger.Error("消息记录反序列化失败", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if smsMessages, exists := smsg[int64(tousr.ID)]; exists {
		for _, msg := range smsMessages {
			resp.Messages = append(resp.Messages, &base.Message{
				Content:    msg[0],
				CreateTime: msg[1],
				UserId:     fmt.Sprintf("%d", usr.ID),
			})
		}
	}
	if amsMessages, exists := amsg[int64(usr.ID)]; exists {
		for _, msg := range amsMessages {
			resp.Messages = append(resp.Messages, &base.Message{
				Content:    msg[0],
				CreateTime: msg[1],
				UserId:     fmt.Sprintf("%d", tousr.ID),
			})
		}
	}
	cacheddatabytes, err := json.Marshal(resp.Messages)
	if err != nil {
		logger.Error("好友列表数据序列化失败：", zap.Error(err))
	}
	err = redis.SetKey(ctx, cacheKey, cacheddatabytes, 10*time.Minute)
	if err != nil {
		logger.Error("缓存设置失败：", zap.Error(err))
	}
	resp.Base = utils.WrapWithSuccess("获取聊天记录成功")
	return resp, nil
}

// DeleteMessage implements the UserServiceImpl interface.
func (s *UserServiceImpl) DeleteMessage(ctx context.Context, req *user.DeleteMessageRequest) (resp *user.DeleteMessageResponse, err error) {
	resp = new(user.DeleteMessageResponse)
	usra, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	usrb, err := db.GetUserById(req.TouserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	err = db.DeleteMessage(usra, usrb)
	if err != nil {
		logger.Error("删除聊天记录失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	cacheKey := fmt.Sprintf("messagelist%d+%d", req.UserId, req.TouserId)
	err = redis.DelKey(ctx, cacheKey) // 删除缓存
	if err != nil {
		logger.Error("删除缓存失败：", zap.Error(err))
	}
	resp.Base = utils.WrapWithSuccess("删除聊天记录成功")
	return resp, nil
}

// GetUsersByName implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUsersByName(ctx context.Context, req *user.GetUsersByNameRequest) (resp *user.GetUsersByNameResponse, err error) {
	resp = new(user.GetUsersByNameResponse)
	resp.Friends = []*base.FriendInfo{}
	users, err := db.GetUsersByName(req.Name)
	if err != nil {
		logger.Error("搜索用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	for _, usr := range users {
		f := new(base.FriendInfo)
		f.Name = usr.Name
		f.UserId = fmt.Sprintf("%d", usr.ID)
		f.Avatar = usr.Avatar
		resp.Friends = append(resp.Friends, f)
	}
	resp.Base = utils.WrapWithSuccess("搜素用户成功")
	return resp, nil
}

// IsFriend implements the UserServiceImpl interface.
func (s *UserServiceImpl) IsFriend(ctx context.Context, req *user.IsFriendRequest) (resp *user.IsFriendResponse, err error) {
	resp = new(user.IsFriendResponse)
	usr, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	ok, err := db.IsFriendJudge(usr, req.TouserId)
	if err != nil {
		logger.Error("查询关系失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Ok = ok
	resp.Base = utils.WrapWithSuccess("获取关系成功")
	return resp, nil
}
