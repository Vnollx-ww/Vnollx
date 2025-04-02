package user

import (
	"Vnollx/pkg/kafka/consumer/user"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

// cd D:
// cd kafka_2.13-3.9.0
// .\bin\windows\kafka-server-start.bat .\config\server.properties
// Kafka 生产者结构体
type KafkaProducer struct {
	producer sarama.SyncProducer
}

// 初始化 Kafka 生产者
func NewKafkaProducer(brokerList []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	return &KafkaProducer{producer: producer}, nil
}
func (k *KafkaProducer) SendSendMessageEvent(UserId int64, Content string, TouserId int64) error {
	message := &user.SendMessageMessage{UserId: UserId, Content: Content, ToUserId: TouserId}
	msgBytes, err := json.Marshal(message)
	msg := &sarama.ProducerMessage{
		Topic: "send-message",
		Value: sarama.StringEncoder(msgBytes),
	}
	_, _, err = k.producer.SendMessage(msg)
	if err != nil {
		log.Printf("发送 发送消息请求 到kafka失败: %v", err)
		return err
	}
	log.Println("发送 发送消息请求 到kafka成功")
	return nil
}
func (k *KafkaProducer) SendFriendApplicationEvent(userId int64, TouserId int64) error {
	message := &user.SendFriendApplicationMessage{UserID: userId, ToUserId: TouserId}
	msgBytes, err := json.Marshal(message)
	msg := &sarama.ProducerMessage{
		Topic: "send-friendapplication",
		Value: sarama.StringEncoder(msgBytes),
	}
	_, _, err = k.producer.SendMessage(msg)
	if err != nil {
		log.Printf("发送 好友请求 到kafka失败: %v", err)
		return err
	}
	log.Println("发送 好友请求 到kafka成功")
	return nil
}
