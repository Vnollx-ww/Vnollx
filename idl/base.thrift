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
    6: string documentcount;
    7: string videocount;
    8: string picturecount;
    9: string musiccount;
    10:string othercount;
    11:string foldercount;
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
struct NotificationInfo{
    1:string title;
    2:string content;
    3:string create_time;
    4:string notification_id;
}
struct FriendInfo{
    1:string name;
    2:string avatar;
    3:string user_id;
}
struct Message{
    1:string content;
    2:string create_time;
    3:string user_id;
}
struct FriendApplication{
    1:string name;
    2:string avatar;
    3:string user_id;
}