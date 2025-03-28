package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

// 定义全局变量
var (
	region          string // 存储区域
	bucketName      string // 存储空间
	accessKeyID     string // 访问密钥ID
	accessKeySecret string // 访问密钥Secret
	objectName      string // 对象名称
	filePath        string // 本地图片路径
)

// init函数用于初始化命令行参数
func init() {
	flag.StringVar(&region, "region", "cn-hangzhou", "The region in which the bucket is located.")
	flag.StringVar(&accessKeyID, "access-key-id", "LTAI5tJ2CXkEVi7qPtXdTdkk", "Your Access Key ID.")
	flag.StringVar(&accessKeySecret, "access-key-secret", "29imopRekuhhdRn0my2hldh46uWdxI", "Your Access Key Secret.")
	flag.StringVar(&objectName, "object", "脚本.txt", "The name of the object in OSS (e.g., image.jpg).")
	flag.StringVar(&bucketName, "bucket", "vnollx", "The bucket name.")
	flag.StringVar(&filePath, "file", "C:\\Users\\吴恩宇\\Desktop\\脚本.txt", "Path to the image file you want to upload.")
}

func main() {
	// 解析命令行参数
	flag.Parse()

	// 检查bucket名称是否为空
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	// 检查region是否为空
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	// 检查Access Key ID 和 Access Key Secret是否为空
	if len(accessKeyID) == 0 || len(accessKeySecret) == 0 {
		log.Fatalf("Access Key ID and Access Key Secret are required")
	}

	// 检查object名称是否为空
	if len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, object name required")
	}

	// 检查文件路径是否为空
	if len(filePath) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, file path required")
	}

	// 打开本地文件（图片）
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	// 创建静态凭证提供者
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

	// 打印上传对象的结果
	log.Printf("put object result:%#v\n", result)
	fmt.Printf("File '%s' uploaded successfully to OSS bucket '%s' as '%s'\n", filePath, bucketName, objectName)
}
