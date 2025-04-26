package global

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	PSNRefreshToken string
	SteamApiKey     string
	RedisUrl        string
	RedisPassword   string
	XboxApiKey      string
}

func Load() *Config {
	viper.AutomaticEnv()

	// 直接设置环境变量前缀（这样就不需要psn.前缀了）
	viper.SetEnvPrefix("") // 设置环境变量前缀为空

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
	}

	// fmt.Println("REFRESH TOKEN HERE: ", viper.GetString("PSN_REFRESH_TOKEN"))

	return &Config{
		PSNRefreshToken: viper.GetString("PSN_REFRESH_TOKEN"),
		SteamApiKey:     viper.GetString("STEAM_API_KEY"),
		XboxApiKey:      viper.GetString("XBOX_API_KEY"),
		RedisUrl:        viper.GetString("REDIS_URL"),
		RedisPassword:   viper.GetString("REDIS_PASSWORD"),
	}
}
