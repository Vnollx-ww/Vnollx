package user

import (
	"Vnollx/dal/db"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"log"
	"time"
)

type SendFriendApplicationConsumer struct {
	consumer sarama.Consumer
	group    string
}
type SendFriendApplicationMessage struct {
	UserID   int64 `json:"user_id"`
	ToUserId int64 `json:"touser_id"`
}

func init() {
	SendFriendApplicationConsumer, err := NewSendFriendApplicationConsumer([]string{kafkaAddr})
	if err != nil {
		log.Println("无法接收消息到 Kafka:", err)
		return
	}
	go func() {
		SendFriendApplicationConsumer.Listen()
	}()
}
func NewSendFriendApplicationConsumer(brokerList []string) (*SendFriendApplicationConsumer, error) {
	consumer, err := sarama.NewConsumer(brokerList, nil)
	if err != nil {
		return nil, err
	}
	return &SendFriendApplicationConsumer{
		consumer: consumer,
		group:    "sendfriendapplication-group",
	}, nil
}
func (c *SendFriendApplicationConsumer) Listen() {
	log.Println("开始监听好友请求")
	config := sarama.NewConfig()
	client, err := sarama.NewClient([]string{kafkaAddr}, config)
	if err != nil {
		log.Fatalf("创建kafka客户端失败: %v", err)
	}
	defer client.Close()
	// 刷新元数据
	err = client.RefreshMetadata("send-friendapplication")
	if err != nil {
		log.Fatalf("刷新元数据失败: %v", err)
	}

	partitionConsumer, err := c.consumer.ConsumePartition("send-friendapplication", 0, sarama.OffsetNewest)
	if err != nil {
		log.Println("friend")
		log.Fatalf("Error consuming partition: %v", err)
	}
	defer partitionConsumer.Close()
	ticker := time.NewTicker(time.Second / 20)
	defer ticker.Stop()
	for msg := range partitionConsumer.Messages() {
		<-ticker.C
		var event SendFriendApplicationMessage
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			logger.Error("反序列化好友申请失败：", zap.Error(err))
			continue
		}
		err = SendFriendApplication(event.UserID, event.ToUserId)
		if err != nil {
			logger.Error("发送好友申请失败：", zap.Error(err))
		}
	}
}
func SendFriendApplication(userID int64, toUserId int64) error {
	log.Println("尝试调用发送好友申请")
	usr, err := db.GetUserById(userID)
	if err != nil {
		return err
	}
	content := fmt.Sprintf("%s(id=%d)请求添加您为好友", usr.Name, usr.ID)
	err = db.CreateNotification("好友请求通知", content, toUserId)
	if err != nil {
		return err
	}
	log.Println("发送好友请求成功！")
	return nil
}
