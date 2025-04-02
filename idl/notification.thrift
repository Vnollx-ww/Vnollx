namespace go notification

include "base.thrift"

struct GetNotificationsByUserIdRequest{
    1:i64 user_id;
}
struct GetNotificationsByUserIdResponse{
    1:base.BaseResponse base;
    2:list<base.NotificationInfo> notifications
}
struct DeleteNotificationRequest{
    1:i64 user_id;
    2:i64 notification_id;
}
struct DeleteNotificationResponse{
    1:base.BaseResponse base;
}
struct SendNotificationRequest{
    1:i64 user_id;
    2:i64 touser_id;
}
struct SendNotificationResponse{
    1:base.BaseResponse base;
}
//kitex -module Vnollx idl/notification.thrift
//kitex -module Vnollx -service Vnollx.notification -use Vnollx/kitex_gen ../../idl/notification.thrift
service NotificationService{
    GetNotificationsByUserIdResponse GetNotificationsByUserId(1:GetNotificationsByUserIdRequest req)
    DeleteNotificationResponse DeleteNotification(1:DeleteNotificationRequest req)
    SendNotificationResponse SendNotification(1:SendNotificationRequest req)
}