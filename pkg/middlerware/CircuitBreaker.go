package middlerware

import (
	base2 "Vnollx/kitex_gen/base"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

var loadedCircuitRules = make(map[string]bool)
var muCircuit sync.Mutex

func init() {
	err := sentinel.InitDefault()
	if err != nil {
		panic(err)
	}
}

func CircuitBreakerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		resource := c.Request.URL.Path

		// 检查该 IP 的熔断规则是否已经加载，如果没有则加载
		if !loadedCircuitRules[resource] {
			muCircuit.Lock()
			if !loadedCircuitRules[resource] {
				rules := []*circuitbreaker.Rule{
					{
						Resource:         resource,
						Strategy:         circuitbreaker.ErrorRatio, // 基于错误率熔断
						Threshold:        0.5,                       // 错误率阈值为 50%
						MinRequestAmount: 50,                        // 最小请求数，达到该数量才开始计算错误率
						StatIntervalMs:   5000,                      // 统计时间窗口为 5 秒
						RetryTimeoutMs:   3000,                      // 熔断后重试时间为 3 秒
					},
				}
				_, err := circuitbreaker.LoadRules(rules)
				if err != nil {
					log.Printf("加载熔断规则IP失败 %s: %v", resource, err)
					c.JSON(http.StatusBadRequest, base2.BaseResponse{
						Code: 400,
						Msg:  "加载熔断规则失败！",
					})
					c.Abort()
					muCircuit.Unlock()
					return
				}
				// 标记该 IP 的规则已加载
				loadedCircuitRules[resource] = true
			}
			muCircuit.Unlock()
		}

		e, b := sentinel.Entry(resource, sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			log.Printf("第 %d 次请求被熔断，IP: %s！", getRequestCount(resource), resource)
			c.JSON(http.StatusServiceUnavailable, base2.BaseResponse{
				Code: 503,
				Msg:  "服务暂时不可用，请稍后重试！",
			})
			c.Abort()
			return
		}
		defer e.Exit()

		c.Next()
		type Response struct {
			Code int `json:"code"`
		}
		var resp Response
		if err := c.ShouldBindJSON(&resp); err == nil && resp.Code == 500 {
			sentinel.TraceError(e, nil)
		}
	}
}

// getRequestCount 辅助函数，用于获取请求计数
func getRequestCount(ip string) int {
	mu.Lock()
	defer mu.Unlock()
	return cnt[ip]
}
