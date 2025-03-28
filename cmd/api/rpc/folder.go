package rpc

import (
	"Vnollx/kitex_gen/folder"
	"Vnollx/kitex_gen/folder/folderservice"
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
	folderClient folderservice.Client
)

func InitFolder(config *viper.Config) {
	etcdAddr := fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	serviceName := config.Viper.GetString("server.name")
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}
	c, err := folderservice.NewClient(
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
	folderClient = c
}
func CreateFolder(ctx context.Context, req *folder.CreateFolderRequest) (*folder.CreateFolderResponse, error) {
	return folderClient.CreateFolder(ctx, req)
}
func DeleteFolder(ctx context.Context, req *folder.DeleteFolderRequest) (*folder.DeleteFolderResponse, error) {
	return folderClient.DeleteFolder(ctx, req)
}
func UpdateFolderInfo(ctx context.Context, req *folder.UpdateFolderInfoRequest) (*folder.UpdateFolderInfoResponse, error) {
	return folderClient.UpdateFolderInfo(ctx, req)
}
func GetFolderInfo(ctx context.Context, req *folder.GetFolderInfoRequest) (*folder.GetFolderInfoResponse, error) {
	return folderClient.GetFolderInfo(ctx, req)
}
func GetFoldersInfo(ctx context.Context, req *folder.GetFoldersInfoRequest) (*folder.GetFoldersInfoResponse, error) {
	return folderClient.GetFoldersInfo(ctx, req)
}
