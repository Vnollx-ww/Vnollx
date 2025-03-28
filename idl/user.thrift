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
struct UserRegisterRequest {
    1: string name;
    2: string phone;
    3: string password;
    4: string avatar;
}
struct UserRegisterResponse {
    1:base.BaseResponse base;
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
    4:string password;
    5:string avatar;
}
struct UpdateUserInfoResponse{
    1:base.BaseResponse base;
}
//kitex -module Vnollx idl/user.thrift
//kitex -module Vnollx -service Vnollx.user -use Vnollx/kitex_gen ../../idl/user.thrift
service UserService{
    UserLoginResponse UserLogin(1:UserLoginRequest req)
    UserRegisterResponse UserResgiter(1:UserRegisterRequest req)
    GetUserInfoResponse GetUserInfo(1:GetUserInfoRequest req)
    UpdateUserInfoResponse UpdateUserInfo(1:UpdateUserInfoRequest req)
}