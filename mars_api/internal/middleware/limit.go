package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/common/pkg/resp"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/store/memory"

	"strings"
	"time"
)

// 如果套用了cdn 替换真实ip header
func getIP(c *gin.Context) string {
	// 从 X-Real-IP 或 X-Forwarded-For 获取客户端的真实 IP
	realIP := c.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// 如果存在 X-Forwarded-For，使用第一个地址
	forwardedFor := c.GetHeader("X-Forwarded-For")
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// 默认使用远程地址
	return c.ClientIP()
}

func RateLimiter() gin.HandlerFunc {
	// 例如每秒 5 次请求
	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  5,
	}

	// 使用内存存储
	store := memory.NewStore()

	// 创建限速器实例
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		ip := getIP(c)
		context, err := instance.Get(c, ip)
		if err != nil {
			resp.Return(c, 500, "Internal Server Error", nil)
			return
		}

		if context.Reached {
			resp.Return(c, 500, "too many requests", nil)
			return
		}

		c.Next()
	}
}
