package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/google/common/pkg/client"
	"github.com/google/common/pkg/config"
	"github.com/google/common/pkg/models"
	"github.com/google/mars_api/internal/conf"
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

	// 初始化数据库
	postgresClient, err := client.PostgresClient(appConfig.PostgresConfiguration, nil)
	if err != nil {
		log.Error().Msg("Failed to connect to postgres")
		panic(err)
	}

	// migration
	err = postgresClient.AutoMigrate(&models.LangPhrase{})
	if err != nil {
		panic(err)
	}
}
