namespace go base
struct BaseResponse {
     1: i32 code;
     2: string msg;
}
struct UserInfo{
    1: string name;
    2: string avatar;
    3: string phone;
    4: i64 id;
    5: string space;
}
struct FileInfo{
    1: string name;
    2: string size;
    3: string postfix;
    4: string parent_folder_id;
    5: string upload_time;
    6: string file_id;
}
struct FolderInfo{
    1:string name;
    2:string parent_folder_id;
    3:string folder_id;
    4:string upload_time;
}