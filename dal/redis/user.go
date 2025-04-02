package redis

import (
	"Vnollx/kitex_gen/base"
	"context"
	"encoding/json"
	"log"
)

func GetUserInfo(ctx context.Context, cachekey string) (*base.UserInfo, error) {
	cacheuserdata, err := client.Get(ctx, cachekey).Result()
	if err == nil && cacheuserdata != "" {
		var cacheduser base.UserInfo
		err := json.Unmarshal([]byte(cacheuserdata), &cacheduser)
		if err != nil {
			log.Println("缓存反序列化失败:", err)
			client.Del(ctx, cachekey)
			return nil, err
		}
		return &cacheduser, nil
	}
	return nil, err
}
func GetFriendList(ctx context.Context, cachekey string) ([]*base.FriendInfo, error) {
	// 尝试从缓存获取数据
	cacheuserdata, err := client.Get(ctx, cachekey).Result()
	if err == nil && cacheuserdata != "" {
		var cacheduserfriend []base.FriendInfo
		err := json.Unmarshal([]byte(cacheuserdata), &cacheduserfriend)
		if err != nil {
			// 如果反序列化失败，删除缓存并记录错误
			log.Println("缓存反序列化失败:", err)
			client.Del(ctx, cachekey)
			return nil, err
		}
		var result []*base.FriendInfo
		for i := range cacheduserfriend {
			result = append(result, &cacheduserfriend[i])
		}
		return result, nil
	}
	return nil, err
}
func GetMessageListByFriend(ctx context.Context, cachekey string) ([]*base.Message, error) {
	// 尝试从缓存获取数据
	cacheuserdata, err := client.Get(ctx, cachekey).Result()
	if err == nil && cacheuserdata != "" {
		var cacheduserfriend []base.Message
		err := json.Unmarshal([]byte(cacheuserdata), &cacheduserfriend)
		if err != nil {
			// 如果反序列化失败，删除缓存并记录错误
			log.Println("缓存反序列化失败:", err)
			client.Del(ctx, cachekey)
			return nil, err
		}
		var result []*base.Message
		for i := range cacheduserfriend {
			result = append(result, &cacheduserfriend[i])
		}
		return result, nil
	}
	return nil, err
}
