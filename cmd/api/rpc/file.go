package rpc

import (
	"Vnollx/kitex_gen/file"
	"Vnollx/kitex_gen/file/fileservice"
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
	fileClient fileservice.Client
)

func InitFile(config *viper.Config) {
	etcdAddr := fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	serviceName := config.Viper.GetString("server.name")
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}
	c, err := fileservice.NewClient(
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
	fileClient = c
}
func UploadFile(ctx context.Context, req *file.UploadFileRequest) (*file.UploadFileResponse, error) {
	return fileClient.UploadFile(ctx, req)
}
func DeleteFile(ctx context.Context, req *file.DeleteFileRequest) (*file.DeleteFileResponse, error) {
	return fileClient.DeleteFile(ctx, req)
}
func UpdateFileInfo(ctx context.Context, req *file.UpdateFileInfoRequest) (*file.UpdateFileInfoResponse, error) {
	return fileClient.UpdateFileInfo(ctx, req)
}
func GetFileInfo(ctx context.Context, req *file.GetFileInfoRequest) (*file.GetFileInfoResponse, error) {
	return fileClient.GetFileInfo(ctx, req)
}
func GetFilesInfo(ctx context.Context, req *file.GetFilesInfoRequest) (*file.GetFilesInfoResponse, error) {
	return fileClient.GetFilesInfo(ctx, req)
}
