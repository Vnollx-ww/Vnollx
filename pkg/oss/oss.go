package main

import (
	"context"
	"flag"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"log"
	"mime/multipart"
)

var (
	region          string // 存储区域
	bucketName      string // 存储空间
	accessKeyID     string // 访问密钥ID
	accessKeySecret string // 访问密钥Secret
	objectName      string // 对象名称
)

func init() {
	flag.StringVar(&region, "region", "cn-hangzhou", "The region in which the bucket is located.")
	flag.StringVar(&accessKeyID, "access-key-id", "LTAI5tJ2CXkEVi7qPtXdTdkk", "Your Access Key ID.")
	flag.StringVar(&accessKeySecret, "access-key-secret", "29imopRekuhhdRn0my2hldh46uWdxI", "Your Access Key Secret.")
	flag.StringVar(&bucketName, "bucket", "vnollx", "The bucket name.")
}

func UploadFile(filename string, file multipart.File) {
	flag.Parse()
	creds := credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")

	// 加载默认配置并设置凭证提供者和区域
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(creds).
		WithRegion(region)

	// 创建OSS客户端
	client := oss.NewClient(cfg)

	// 创建上传对象的请求
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(bucketName), // 存储空间名称
		Key:    oss.Ptr(objectName), // 对象名称
		Body:   file,                // 上传的文件内容
	}

	// 执行上传对象的请求
	result, err := client.PutObject(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put object %v", err)
	}
	log.Printf("put object result:%#v\n", result)
}
