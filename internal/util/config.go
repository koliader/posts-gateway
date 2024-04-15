package util

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Config struct {
	GinMode         string `mapstructure:"GIN_MODE"`
	AuthGrpcService string `mapstructure:"AUTH_GRPC_SERVICE"`
	ServerAddress   string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	if config.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	return
}
