package rpc

import (
	"Vnollx/kitex_gen/user"
	"Vnollx/kitex_gen/user/userservice"
	"Vnollx/pkg/viper"
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	etcd "github.com/kitex-contrib/registry-etcd"
	"time"
)

var (
	userClient userservice.Client
)

func InitUser(config *viper.Config) {
	etcdAddr := fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	serviceName := config.Viper.GetString("server.name")
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}
	c, err := userservice.NewClient(
		serviceName,
		client.WithRPCTimeout(30*time.Second),             // rpc timeout
		client.WithConnectTimeout(30000*time.Millisecond), // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver

		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}
func Register(ctx context.Context, req *user.UserRegisterRequest) (*user.UserRegisterResponse, error) {
	return userClient.UserResgiter(ctx, req)
}
func Login(ctx context.Context, req *user.UserLoginRequest) (*user.UserLoginResponse, error) {
	return userClient.UserLogin(ctx, req)
}
func GetUserInfo(ctx context.Context, req *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	return userClient.GetUserInfo(ctx, req)
}
func UpdateUserInfo(ctx context.Context, req *user.UpdateUserInfoRequest) (*user.UpdateUserInfoResponse, error) {
	return userClient.UpdateUserInfo(ctx, req)
}
func UpdatePassword(ctx context.Context, req *user.UpdatePasswordRequest) (*user.UpdatePasswordResponse, error) {
	return userClient.UpdatePassword(ctx, req)
}
func LoginByCode(ctx context.Context, req *user.UserLoginByCodeRequest) (*user.UserLoginByCodeResponse, error) {
	return userClient.UserLoginByCode(ctx, req)
}
func GenerateCaptcha(ctx context.Context, req *user.GenerateCaptchaRequest) (*user.GenerateCaptchaResponse, error) {
	return userClient.GenerateCaptcha(ctx, req)
}
func AddFriend(ctx context.Context, req *user.AddFriendRequest) (*user.AddFriendResponse, error) {
	return userClient.AddFriend(ctx, req)
}
func DeleteFriend(ctx context.Context, req *user.DeleteFriendRequest) (*user.DeleteFriendResponse, error) {
	return userClient.DeleteFriend(ctx, req)
}
func SendMessage(ctx context.Context, req *user.SendMessageRequest) (*user.SendMessageResponse, error) {
	return userClient.SendMessage(ctx, req)
}
func GetFriendList(ctx context.Context, req *user.GetFriendListRequest) (*user.GetFriendListResponse, error) {
	return userClient.GetFriendList(ctx, req)
}
func GetMessageList(ctx context.Context, req *user.GetMessageListRequest) (*user.GetMessageListResponse, error) {
	return userClient.GetMessageList(ctx, req)
}
func DeleteMessage(ctx context.Context, req *user.DeleteMessageRequest) (*user.DeleteMessageResponse, error) {
	return userClient.DeleteMessage(ctx, req)
}
func GetUsersByName(ctx context.Context, req *user.GetUsersByNameRequest) (*user.GetUsersByNameResponse, error) {
	return userClient.GetUsersByName(ctx, req)
}
func IsFriend(ctx context.Context, req *user.IsFriendRequest) (*user.IsFriendResponse, error) {
	return userClient.IsFriend(ctx, req)
}
