package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/common/pkg/open_telemetry"
)

func (a *ApiServer) Router() {
	a.app.GET("/health", a.HealthCheck)

	// verify 验证码 短信 等
	{
		a.app.GET("/captcha", a.CaptchaImage)
	}
}

func (a *ApiServer) HealthCheck(ctx *gin.Context) {
	_, span := open_telemetry.Tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	ctx.JSON(200, gin.H{
		"message": "ok",
	})
}
