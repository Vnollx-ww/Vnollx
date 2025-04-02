package main

import (
	"Vnollx/kitex_gen/notification/notificationservice"
	"Vnollx/pkg/viper"
	"Vnollx/pkg/zaplog"
	"fmt"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
	"log"
	"net"
)

var (
	config      = viper.Init("notification")
	serviceName = config.Viper.GetString("server.name")
	serviceAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
	etcdAddr    = fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	kafkaAddr   = fmt.Sprintf("%s:%d", config.Viper.GetString("kafka.host"), config.Viper.GetInt("kafka.port"))
	logger      = zaplog.GetLogger()
)

func main() {
	r, err := etcd.NewEtcdRegistry([]string{etcdAddr})
	if err != nil {
		log.Fatal(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", serviceAddr)
	if err != nil {
		log.Fatalln(err)
	}

	// 初始化etcd
	s := notificationservice.NewServer(new(NotificationServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)
	if err := s.Run(); err != nil {
		logger.Fatal("Service stopped with error", zap.String("serviceName", serviceName), zap.String("error", err.Error()))
	}
}
