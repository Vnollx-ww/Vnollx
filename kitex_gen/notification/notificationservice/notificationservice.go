// Code generated by Kitex v0.7.3. DO NOT EDIT.

package notificationservice

import (
	notification "Vnollx/kitex_gen/notification"
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return notificationServiceServiceInfo
}

var notificationServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "NotificationService"
	handlerType := (*notification.NotificationService)(nil)
	methods := map[string]kitex.MethodInfo{
		"GetNotificationsByUserId": kitex.NewMethodInfo(getNotificationsByUserIdHandler, newNotificationServiceGetNotificationsByUserIdArgs, newNotificationServiceGetNotificationsByUserIdResult, false),
		"DeleteNotification":       kitex.NewMethodInfo(deleteNotificationHandler, newNotificationServiceDeleteNotificationArgs, newNotificationServiceDeleteNotificationResult, false),
		"SendNotification":         kitex.NewMethodInfo(sendNotificationHandler, newNotificationServiceSendNotificationArgs, newNotificationServiceSendNotificationResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "notification",
		"ServiceFilePath": `idl\notification.thrift`,
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.7.3",
		Extra:           extra,
	}
	return svcInfo
}

func getNotificationsByUserIdHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*notification.NotificationServiceGetNotificationsByUserIdArgs)
	realResult := result.(*notification.NotificationServiceGetNotificationsByUserIdResult)
	success, err := handler.(notification.NotificationService).GetNotificationsByUserId(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newNotificationServiceGetNotificationsByUserIdArgs() interface{} {
	return notification.NewNotificationServiceGetNotificationsByUserIdArgs()
}

func newNotificationServiceGetNotificationsByUserIdResult() interface{} {
	return notification.NewNotificationServiceGetNotificationsByUserIdResult()
}

func deleteNotificationHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*notification.NotificationServiceDeleteNotificationArgs)
	realResult := result.(*notification.NotificationServiceDeleteNotificationResult)
	success, err := handler.(notification.NotificationService).DeleteNotification(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newNotificationServiceDeleteNotificationArgs() interface{} {
	return notification.NewNotificationServiceDeleteNotificationArgs()
}

func newNotificationServiceDeleteNotificationResult() interface{} {
	return notification.NewNotificationServiceDeleteNotificationResult()
}

func sendNotificationHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*notification.NotificationServiceSendNotificationArgs)
	realResult := result.(*notification.NotificationServiceSendNotificationResult)
	success, err := handler.(notification.NotificationService).SendNotification(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newNotificationServiceSendNotificationArgs() interface{} {
	return notification.NewNotificationServiceSendNotificationArgs()
}

func newNotificationServiceSendNotificationResult() interface{} {
	return notification.NewNotificationServiceSendNotificationResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) GetNotificationsByUserId(ctx context.Context, req *notification.GetNotificationsByUserIdRequest) (r *notification.GetNotificationsByUserIdResponse, err error) {
	var _args notification.NotificationServiceGetNotificationsByUserIdArgs
	_args.Req = req
	var _result notification.NotificationServiceGetNotificationsByUserIdResult
	if err = p.c.Call(ctx, "GetNotificationsByUserId", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DeleteNotification(ctx context.Context, req *notification.DeleteNotificationRequest) (r *notification.DeleteNotificationResponse, err error) {
	var _args notification.NotificationServiceDeleteNotificationArgs
	_args.Req = req
	var _result notification.NotificationServiceDeleteNotificationResult
	if err = p.c.Call(ctx, "DeleteNotification", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) SendNotification(ctx context.Context, req *notification.SendNotificationRequest) (r *notification.SendNotificationResponse, err error) {
	var _args notification.NotificationServiceSendNotificationArgs
	_args.Req = req
	var _result notification.NotificationServiceSendNotificationResult
	if err = p.c.Call(ctx, "SendNotification", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
