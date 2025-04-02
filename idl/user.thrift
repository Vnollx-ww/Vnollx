namespace go user

include "base.thrift"

struct UserLoginRequest {
    1: string phone;
    2: string password;
}
struct UserLoginResponse {
    1:base.BaseResponse base;
    2:string token;
}
struct UserLoginByCodeRequest{
    1: string phone;
    2: string captcha;
}
struct UserLoginByCodeResponse{
    1:base.BaseResponse base;
    2:string token;
}
struct UserRegisterRequest {
    1: string name;
    2: string phone;
    3: string password;
    4: string postfix;
}
struct UserRegisterResponse {
    1:base.BaseResponse base;
    2:string user_id;
}
struct GetUserInfoRequest{
    1:i64 user_id;
}
struct GetUserInfoResponse{
    1:base.BaseResponse base;
    2:base.UserInfo userinfo;
}
struct UpdateUserInfoRequest{
    1:i64 user_id;
    2:string name;
    3:string phone;
}
struct UpdateUserInfoResponse{
    1:base.BaseResponse base;
}
struct UpdatePasswordRequest{
    1:i64 user_id;
    2:string old_password;
    3:string new_password;
}
struct UpdatePasswordResponse{
    1:base.BaseResponse base;
}
struct GenerateCaptchaRequest{
    1:string phone;
}
struct GenerateCaptchaResponse{
    1:base.BaseResponse base;
}
struct SendMessageRequest{
    1:string content;
    2:i64 user_id;
    3:i64 touser_id;
}
struct SendMessageResponse{
    1:base.BaseResponse base;
}
struct GetFriendListRequest{
    1:i64 user_id;
}
struct GetFriendListResponse{
    1:base.BaseResponse base;
    2:list<base.FriendInfo> friends;
}
struct AddFriendRequest{
    1:i64 user_id;
    2:i64 touser_id;
}
struct AddFriendResponse{
    1:base.BaseResponse base;
}
struct DeleteFriendRequest{
    1:i64 user_id;
    2:i64 touser_id;
}
struct DeleteFriendResponse{
    1:base.BaseResponse base;
}
struct GetMessageListRequest{
    1:i64 user_id;
    2:i64 touser_id;
}
struct GetMessageListResponse{
    1:base.BaseResponse base;
    2:list<base.Message> messages;
}
struct DeleteMessageRequest{
    1:i64 user_id;
    2:i64 touser_id;
}
struct DeleteMessageResponse{
    1:base.BaseResponse base;
}
struct GetUsersByNameRequest{
    1:string name;
}
struct GetUsersByNameResponse{
    1:base.BaseResponse base;
    2:list<base.FriendInfo> friends;
}
struct IsFriendRequest{
    1:i64 user_id;
    2:i64 touser_id;
}
struct IsFriendResponse{
    1:base.BaseResponse base;
    2:bool ok;
}
//kitex -module Vnollx idl/user.thrift
//kitex -module Vnollx -service Vnollx.user -use Vnollx/kitex_gen ../../idl/user.thrift
service UserService{
    UserLoginResponse UserLogin(1:UserLoginRequest req)
    UserLoginByCodeResponse UserLoginByCode(1:UserLoginByCodeRequest req)
    UserRegisterResponse UserResgiter(1:UserRegisterRequest req)
    GetUserInfoResponse GetUserInfo(1:GetUserInfoRequest req)
    UpdateUserInfoResponse UpdateUserInfo(1:UpdateUserInfoRequest req)
    UpdatePasswordResponse UpdatePassword(1:UpdatePasswordRequest req)
    GenerateCaptchaResponse GenerateCaptcha(1:GenerateCaptchaRequest req)
    //好友功能
    AddFriendResponse AddFriend(1:AddFriendRequest req)
    DeleteFriendResponse DeleteFriend(1:DeleteFriendRequest req)
    SendMessageResponse SendMessage(1: SendMessageRequest req)
    GetFriendListResponse GetFriendList(1:GetFriendListRequest req)
    GetMessageListResponse GetMessageList(1:GetMessageListRequest req)
    DeleteMessageResponse DeleteMessage(1:DeleteMessageRequest req)
    GetUsersByNameResponse GetUsersByName(1:GetUsersByNameRequest req)
    IsFriendResponse IsFriend(1:IsFriendRequest req)
}