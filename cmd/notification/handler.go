package main

import (
	"Vnollx/dal/db"
	"Vnollx/kitex_gen/base"
	notification "Vnollx/kitex_gen/notification"
	user2 "Vnollx/pkg/kafka/producer"
	"Vnollx/utils"
	"context"
	"fmt"
	"go.uber.org/zap"
)

// NotifictionServiceImpl implements the last service interface defined in the IDL.
type NotificationServiceImpl struct{}

// GetNotificationsByUserId implements the NotifictionServiceImpl interface.
func (s *NotificationServiceImpl) GetNotificationsByUserId(ctx context.Context, req *notification.GetNotificationsByUserIdRequest) (resp *notification.GetNotificationsByUserIdResponse, err error) {
	resp = new(notification.GetNotificationsByUserIdResponse)
	resp.Notifications = []*base.NotificationInfo{}
	notifications, err := db.GetNotificationsByUserId(req.UserId)
	if err != nil {
		logger.Error("获取通知列表失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	for _, no := range notifications {
		n := new(base.NotificationInfo)
		n.Title = no.Title
		n.Content = no.Content
		n.CreateTime = no.CreateTime
		n.NotificationId = fmt.Sprintf("%d", no.ID)
		resp.Notifications = append(resp.Notifications, n)
	}
	resp.Base = utils.WrapWithSuccess("获取通知列表成功")
	return resp, nil
}

// DeleteNotification implements the NotifictionServiceImpl interface.
func (s *NotificationServiceImpl) DeleteNotification(ctx context.Context, req *notification.DeleteNotificationRequest) (resp *notification.DeleteNotificationResponse, err error) {
	resp = new(notification.DeleteNotificationResponse)
	n, err := db.GetNotificationById(req.NotificationId)
	if err != nil {
		logger.Error("获取通知失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	if n == nil {
		logger.Error("通知不存在或已被删除：", zap.Error(err))
		resp.Base = utils.WrapWithLogicError("通知不存在或已被删除")
		return resp, nil
	}
	err = db.DeleteNotification(n)
	if err != nil {
		logger.Error("删除通知失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("删除通知成功")
	return resp, nil
}

// SendNotification implements the NotificationServiceImpl interface.
func (s *NotificationServiceImpl) SendNotification(ctx context.Context, req *notification.SendNotificationRequest) (resp *notification.SendNotificationResponse, err error) {
	resp = new(notification.SendNotificationResponse)
	kafkaProducer, err1 := user2.NewKafkaProducer([]string{kafkaAddr}) //初始化kafka生产者
	var err2 error
	if err1 == nil {
		err2 = kafkaProducer.SendFriendApplicationEvent(req.UserId, req.TouserId)
		if err2 == nil {
			resp.Base = utils.WrapWithSuccess("发送好友请求通知成功")
			return resp, nil
		}
	}
	if err1 != nil {
		logger.Error("kafka生产者创建失败：", zap.Error(err))
	}
	if err2 != nil {
		logger.Error("发送好友请求的请求创建成功，但请求发送到kafka失败：", zap.Error(err))
	}
	usr, err := db.GetUserById(req.UserId)
	if err != nil {
		logger.Error("查询用户失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	content := fmt.Sprintf("%s(id=%d)请求添加您为好友", usr.Name, usr.ID)
	err = db.CreateNotification("好友请求通知", content, req.TouserId)
	if err != nil {
		logger.Error("发送好友请求通知失败：", zap.Error(err))
		resp.Base = utils.WrapWithSystemError("服务器内部错误，请联系管理员")
		return resp, nil
	}
	resp.Base = utils.WrapWithSuccess("发送好友请求通知成功")
	return resp, nil
}
