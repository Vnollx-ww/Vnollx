package rpc

import (
	"Vnollx/kitex_gen/notification"
	"Vnollx/kitex_gen/notification/notificationservice"
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
	notificationClient notificationservice.Client
)

func InitNotification(config *viper.Config) {
	etcdAddr := fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	serviceName := config.Viper.GetString("server.name")
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}
	c, err := notificationservice.NewClient(
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
	notificationClient = c
}
func GetNotificationsByUserId(ctx context.Context, req *notification.GetNotificationsByUserIdRequest) (*notification.GetNotificationsByUserIdResponse, error) {
	return notificationClient.GetNotificationsByUserId(ctx, req)
}
func DeleteNotification(ctx context.Context, req *notification.DeleteNotificationRequest) (*notification.DeleteNotificationResponse, error) {
	return notificationClient.DeleteNotification(ctx, req)
}
func SendNotification(ctx context.Context, req *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error) {
	return notificationClient.SendNotification(ctx, req)
}
