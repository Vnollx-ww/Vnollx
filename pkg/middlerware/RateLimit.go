package middlerware

import (
	base2 "Vnollx/kitex_gen/base"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

var loadedIPRules = make(map[string]bool)
var cnt = make(map[string]int)
var mu sync.Mutex

func init() {
	err := sentinel.InitDefault()
	if err != nil {
		panic(err)
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//clientIP := c.ClientIP()
		resource := c.Request.URL.Path

		// 加锁保护对 cnt 的并发访问，确保每次请求计数正确增加
		mu.Lock()
		cnt[resource]++
		//currentCount := cnt[resource]
		mu.Unlock()

		// 检查该 IP 的规则是否已经加载，如果没有则加载
		if !loadedIPRules[resource] {
			mu.Lock()
			if !loadedIPRules[resource] {
				_, err := flow.LoadRules([]*flow.Rule{
					{
						Resource:               resource,
						TokenCalculateStrategy: flow.Direct,
						ControlBehavior:        flow.Reject,
						Threshold:              1000, // 每秒允许通过的请求数为 10
						StatIntervalInMs:       1000,
					},
				})
				if err != nil {
					log.Printf("加载限流规则接口失败 %s: %v", resource, err)
					c.JSON(http.StatusBadRequest, base2.BaseResponse{
						Code: 400,
						Msg:  "加载限流规则失败！",
					})
					c.Abort()
					mu.Unlock()
					return
				}
				// 标记该 IP 的规则已加载
				loadedIPRules[resource] = true
			}
			mu.Unlock()
		}

		e, b := sentinel.Entry(resource, sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			//log.Printf("第 %d 次请求繁忙，请稍后重试，接口: %s！", currentCount, resource)
			c.JSON(http.StatusBadRequest, base2.BaseResponse{
				Code: 400,
				Msg:  "请求繁忙，请稍后重试！",
			})
			c.Abort()
			return
		}
		c.Next()
		e.Exit()
	}
}
