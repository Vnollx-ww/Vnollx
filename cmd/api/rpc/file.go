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
func GetAllFiles(ctx context.Context, req *file.GetALLFilesRequest) (*file.GetAllFilesResponse, error) {
	return fileClient.GetAllFiles(ctx, req)
}
func GetVideoFiles(ctx context.Context, req *file.GetVideoFilesRequest) (*file.GetVideoFilesResponse, error) {
	return fileClient.GetVideoFiles(ctx, req)
}
func GetMusicFiles(ctx context.Context, req *file.GetMusicFilesRequest) (*file.GetMusicFilesResponse, error) {
	return fileClient.GetMusicFiles(ctx, req)
}
func GetPictureFiles(ctx context.Context, req *file.GetPictureFilesRequest) (*file.GetPictureFilesResponse, error) {
	return fileClient.GetPictureFiles(ctx, req)
}
func GetDocumentFiles(ctx context.Context, req *file.GetDocumentFilesRequest) (*file.GetDocumentFilesResponse, error) {
	return fileClient.GetDocumentFiles(ctx, req)
}
func GetOtherFiles(ctx context.Context, req *file.GetOtherFilesRequest) (*file.GetOtherFilesResponse, error) {
	return fileClient.GetOtherFiles(ctx, req)
}
func GetFilesByName(ctx context.Context, req *file.GetFilesByNameRequest) (*file.GetFilesByNameResponse, error) {
	return fileClient.GetFilesByName(ctx, req)
}
func CreateShare(ctx context.Context, req *file.CreateShareRequest) (*file.CreateShareResponse, error) {
	return fileClient.CreateShare(ctx, req)
}
func GetFileByCode(ctx context.Context, req *file.GetFileByCodeRequest) (*file.GetFileByCodeResponse, error) {
	return fileClient.GetFileByCode(ctx, req)
}
