namespace go folder

include "base.thrift"

struct CreateFolderRequest{
    1:i64 user_id
    2:string name;
    3:i64 parent_folder_id;
}
struct CreateFolderResponse{
    1:base.BaseResponse base;
}
struct DeleteFolderRequest{
    1:i64 user_id;
    2:i64 folder_id;
}
struct DeleteFolderResponse{
    1:base.BaseResponse base;
}
struct UpdateFolderInfoRequest{
    1:i64 user_id;
    2:i64 folder_id;
    3:string name;
    4:i64 parent_folder_id;
}
struct UpdateFolderInfoResponse{
    1:base.BaseResponse base;
}
struct GetFolderInfoRequest{
    1:i64 folder_id;
}
struct GetFolderInfoResponse{
    1:base.BaseResponse base;
    2:base.FolderInfo folder;
}
struct GetFoldersInfoRequest{
    1:i64 parent_folder_id;
    2:i64 user_id;
}
struct GetFoldersInfoResponse{
    1:base.BaseResponse base;
    2:list<base.FolderInfo> folders;
}
//kitex -module Vnollx idl/folder.thrift
//kitex -module Vnollx -service Vnollx.folder -use Vnollx/kitex_gen ../../idl/folder.thrift
service FolderService{
    CreateFolderResponse CreateFolder(1:CreateFolderRequest req)
    DeleteFolderResponse DeleteFolder(1:DeleteFolderRequest req)
    UpdateFolderInfoResponse UpdateFolderInfo(1:UpdateFolderInfoRequest req)
    GetFolderInfoResponse GetFolderInfo(1:GetFolderInfoRequest req)
    GetFoldersInfoResponse GetFoldersInfo(1:GetFoldersInfoRequest req)
}