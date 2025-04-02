package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"log"
)

// D:\minio.exe server D:\minio --console-address :9090
var MinioClient *minio.Client

func init() {
	//endpoint := "localhost:9000"
	//endpoint := "my-minio:9000"
	endpoint := "106.54.223.38:9000"
	accessKeyID := "vnollxvnollx"
	secretAccessKey := "vnollxvnollxvnollx"

	// 初始化一个minio客户端对象
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Fatal("初始化MinioClient错误:", zap.Error(err))
	}
	MinioClient = minioClient
}
func DeleteFileMinio(c context.Context, file_id int64, postfix string) error {
	bucketName := "file"
	objectName := fmt.Sprintf("FileId=%d.%s", file_id, postfix)
	// 删除文件
	err := MinioClient.RemoveObject(c, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Fatalf("删除文件失败: %v", err)
		return err
	}
	return nil
}
