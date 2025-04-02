namespace go file

include "base.thrift"
/*struct DownloadFileRequest{
1:i64 user_id
2:i64 file_id
}
struct DownloadFileResponse{
1: i32    status_code
2: string status_msg
3: bool success
}*/



struct UploadFileRequest{
    1:i64 user_id;
    2:i64 parent_folder_id;
    3:string name;
    4:string size;
    5:string source;
    6:string postfix;
}
struct UploadFileResponse{
    1:base.BaseResponse base;
    2:string file_id;
}
struct DeleteFileRequest{
    1:i64 user_id;
    2:i64 file_id;
}
struct DeleteFileResponse{
    1:base.BaseResponse base;
}
struct UpdateFileInfoRequest{
    1:i64 user_id;
    2:i64 file_id;
    3:string name;
    4:i64 parent_folder_id;
}
struct UpdateFileInfoResponse{
    1:base.BaseResponse base;
}
struct GetFileInfoRequest{
    1:i64 file_id;
}
struct GetFileInfoResponse{
    1:base.BaseResponse base;
    2:base.FileInfo file;
}
struct GetFilesInfoRequest{
    1:i64 parent_folder_id;
    2:i64 user_id;
}
struct GetFilesInfoResponse{
    1:base.BaseResponse base;
    2:list<base.FileInfo> files;
}
struct GetALLFilesRequest{
    1:i64 user_id;
    2:string name;
}
struct GetAllFilesResponse{
    1:base.BaseResponse base;
    2:list<base.FileInfo> files;
}
struct GetVideoFilesRequest{
    1:i64 user_id;
    2:string name;
}
struct GetVideoFilesResponse{
    1:base.BaseResponse base;
    2:list<base.FileInfo>files;
}
struct GetMusicFilesRequest{
    1:i64 user_id;
    2:string name;
}
struct GetMusicFilesResponse{
    1:base.BaseResponse base;
    2:list<base.FileInfo>files;
}
struct GetPictureFilesRequest{
    1:i64 user_id;
    2:string name;
}
struct GetPictureFilesResponse{
    1:base.BaseResponse base;
    2:list<base.FileInfo>files;
}
struct GetDocumentFilesRequest{
    1:i64 user_id;
    2:string name;
}
struct GetDocumentFilesResponse{
    1:base.BaseResponse base;
    2:list<base.FileInfo>files;
}
struct GetOtherFilesRequest{
    1:i64 user_id;
    2:string name;
}
struct GetOtherFilesResponse{
    1:base.BaseResponse base;
    2:list<base.FileInfo>files;
}
struct GetFilesByNameRequest{
    1:string name;
    2:i64 user_id;
}
struct GetFilesByNameResponse{
    1:base.BaseResponse base;
    2:list<base.FileInfo> files;
}
struct CreateShareRequest{
    1:i64 user_id;
    2:i64 file_id;
}
struct CreateShareResponse{
    1:base.BaseResponse base;
    2:string code;
}
struct GetFileByCodeRequest{
    1:string code;
}
struct GetFileByCodeResponse{
    1:base.BaseResponse base;
    2:base.FileInfo file;
}
//kitex -module Vnollx idl/file.thrift
//kitex -module Vnollx -service Vnollx.file -use Vnollx/kitex_gen ../../idl/file.thrift
service FileService{
    UploadFileResponse UploadFile(1:UploadFileRequest req)
    DeleteFileResponse DeleteFile(1:DeleteFileRequest req)
    UpdateFileInfoResponse UpdateFileInfo(1:UpdateFileInfoRequest req)
    GetFileInfoResponse GetFileInfo(1:GetFileInfoRequest req)
    GetFilesInfoResponse GetFilesInfo(1:GetFilesInfoRequest req)
    GetAllFilesResponse GetAllFiles(1:GetALLFilesRequest req)
    GetVideoFilesResponse GetVideoFiles(1:GetVideoFilesRequest req)
    GetMusicFilesResponse GetMusicFiles(1:GetMusicFilesRequest req)
    GetPictureFilesResponse GetPictureFiles(1:GetPictureFilesRequest req)
    GetDocumentFilesResponse GetDocumentFiles(1:GetDocumentFilesRequest req)
    GetOtherFilesResponse GetOtherFiles(1:GetOtherFilesRequest req)
    GetFilesByNameResponse GetFilesByName(1:GetFilesByNameRequest req)
    CreateShareResponse CreateShare(1:CreateShareRequest req)
    GetFileByCodeResponse GetFileByCode(1:GetFileByCodeRequest req)
}