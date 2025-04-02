package user

import (
	"Vnollx/dal/db"
	"Vnollx/pkg/viper"
	"Vnollx/pkg/zaplog"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"log"
	"time"
)

/*
PS D:\kafka_2.12-3.9.0> .\bin\windows\kafka-topics.bat --delete --topic send-message --bootstrap-server localhost:9092
PS D:\kafka_2.12-3.9.0> .\bin\windows\kafka-topics.bat --create --topic send-message --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
*/
var (
	config    = viper.Init("user")
	kafkaAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("kafka.host"), config.Viper.GetInt("kafka.port"))
)
var logger *zap.Logger = zaplog.GetLogger()

type SendMessageConsumer struct {
	consumer sarama.Consumer
	group    string
}
type SendMessageMessage struct {
	UserId   int64  `json:"user_id"`
	Content  string `json:"content"`
	ToUserId int64  `json:"touser_id"`
}

func init() {
	SendMessageConsumer, err := NewSendMessageConsumer([]string{kafkaAddr})
	if err != nil {
		log.Println("无法接收消息到 Kafka:", err)
		return
	}
	go func() {
		SendMessageConsumer.Listen()
	}()
}
func NewSendMessageConsumer(brokerList []string) (*SendMessageConsumer, error) {
	consumer, err := sarama.NewConsumer(brokerList, nil)
	if err != nil {
		return nil, err
	}
	return &SendMessageConsumer{
		consumer: consumer,
		group:    "sendmessage-group",
	}, nil
}
func (c *SendMessageConsumer) Listen() {
	log.Println("开始监听发送消息请求")
	config := sarama.NewConfig()
	client, err := sarama.NewClient([]string{kafkaAddr}, config)
	if err != nil {
		log.Fatalf("创建kafka客户端失败: %v", err)
	}
	defer client.Close()
	// 刷新元数据
	err = client.RefreshMetadata("send-message")
	if err != nil {
		log.Fatalf("刷新元数据失败: %v", err)
	}
	partitionConsumer, err := c.consumer.ConsumePartition("send-message", 0, sarama.OffsetNewest)
	if err != nil {
		log.Println("message")
		log.Fatalf("Error consuming partition: %v", err)
	}
	defer partitionConsumer.Close()
	ticker := time.NewTicker(time.Second / 100)
	defer ticker.Stop()
	for msg := range partitionConsumer.Messages() {
		<-ticker.C
		var event SendMessageMessage
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			logger.Error("反序列化发送消息事件失败：", zap.Error(err))
			continue
		}
		err = SendMessage(event.UserId, event.Content, event.ToUserId)
		if err != nil {
			logger.Error("发送消息失败：", zap.Error(err))
		}
	}
}

func SendMessage(UserId int64, content string, ToUserId int64) error {
	log.Println("尝试调用发送消息")
	usr, err := db.GetUserById(UserId)
	if err != nil {
		return err
	}
	tousr, err := db.GetUserById(ToUserId)
	if err != nil {
		return err
	}
	err = db.SendMessage(usr, tousr, content)
	if err != nil {
		return err
	}
	return nil
}
