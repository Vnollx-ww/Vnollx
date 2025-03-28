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
