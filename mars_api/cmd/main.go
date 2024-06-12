package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/common/pkg/verification"
	"strings"

	"github.com/google/common/pkg/client"
	"github.com/google/common/pkg/config"
	"github.com/google/common/pkg/logs"
	"github.com/google/common/pkg/open_telemetry"
	"github.com/google/mars_api/internal/conf"
	"github.com/google/mars_api/internal/server"
	"github.com/google/mars_api/internal/storage"
	"github.com/rs/zerolog/log"
)

var configFilename string
var configDirs string

func init() {
	const (
		defaultConfigFilename = "dev_config"
		configUsage           = "Name of the config file, without extension"
		defaultConfigDirs     = "./,./configs/"
		configDirUsage        = "Directories to search for config file, separated by ','"
	)
	flag.StringVar(&configFilename, "c", defaultConfigFilename, configUsage)
	flag.StringVar(&configFilename, "dev_config", defaultConfigFilename, configUsage)
	flag.StringVar(&configDirs, "cPath", defaultConfigDirs, configDirUsage)
}

func main() {
	flag.Parse()

	// config
	var appConfig conf.Config
	err := config.InitConfiguration(configFilename, strings.Split(configDirs, ","), &appConfig)
	if err != nil {
		panic(err)
	}
	indent, err := json.MarshalIndent(appConfig, "", "  ")
	if err == nil {
		fmt.Println(string(indent))
	}
	fmt.Println("Config loaded successfully!")
	// 基础依赖初始化
	// 初始化日志
	logs.InitLog(appConfig.LoggerConfiguration)
	// 初始化链路追踪
	go open_telemetry.InitLog(appConfig.OpenTelemetryConfiguration.Logs)
	openTelemetry := open_telemetry.InitTracerHTTP(appConfig.OpenTelemetryConfiguration.Traces)
	defer func() {
		if err := openTelemetry.Shutdown(context.Background()); err != nil {
			log.Error().Msgf("Failed to connect to postgres %s", err)
		}
	}()
	log.Info().Msg("OpenTelemetry initialized")
	// 初始化数据库
	postgresClient, err := client.PostgresClient(appConfig.PostgresConfiguration, nil)
	if err != nil {
		log.Error().Msg("Failed to connect to postgres")
		panic(err)
	}
	// 初始化缓存
	redisClient, err := client.RedisClient(appConfig.RedisConfiguration)
	if err != nil {
		log.Error().Msg("Failed to connect to redis")
		panic(err)
	}
	storage := storage.NewStorage(redisClient, postgresClient)
	log.Info().Msg("Storage initialized")

	// 初始化验证码服务
	verification.InitVerification(redisClient)
	log.Info().Msg("Verification service initialized")

	// 启动服务
	ser := server.NewServer(storage, appConfig)
	if err := ser.Run(); err != nil {
		log.Error().Msgf("Failed to start server %s", err)
	}
}
