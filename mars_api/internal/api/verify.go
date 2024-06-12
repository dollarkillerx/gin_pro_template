package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/common/pkg/open_telemetry"
	"github.com/google/common/pkg/resp"
	"github.com/google/common/pkg/verification"
	"github.com/rs/zerolog/log"
)

func (a *ApiServer) CaptchaImage(ctx *gin.Context) {
	_, span := open_telemetry.Tracer.Start(ctx, ctx.Request.URL.Path)
	defer span.End()

	id, images, err := verification.Generate()
	if err != nil {
		log.Error().Msgf("get_captcha_failed: %v", err)
		resp.Return(ctx, resp.UnprocessableEntityCode, "get_captcha_failed", nil)
		return
	}

	resp.Return(ctx, resp.SuccessCode, "success", gin.H{
		"cap_id":     id,
		"img_base64": images,
	})
}
