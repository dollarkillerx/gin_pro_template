package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/mars_api/internal/conf"
	"github.com/google/mars_api/internal/middleware"
	"github.com/google/mars_api/internal/storage"
)

type ApiServer struct {
	storage *storage.Storage
	conf    conf.Config
	app     *gin.Engine
}

func NewApiServer(storage *storage.Storage, conf conf.Config) *ApiServer {
	return &ApiServer{storage: storage, conf: conf}
}

func (a *ApiServer) Run() error {

	a.app = gin.New()
	a.app.Use(middleware.HttpRecover())
	a.app.Use(middleware.RateLimiter())
	a.app.Use(gin.Logger())
	a.app.Use(middleware.Cors())

	a.Router()

	return a.app.Run(fmt.Sprintf("127.0.0.1:%s", a.conf.ServiceConfiguration.Port))
}
