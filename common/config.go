package common

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/getsentry/raven-go"
	"github.com/spf13/viper"
)

// SetConfig 配置设置
func SetConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("conf")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("FEIYAN")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Println("Fatal error config file:", err)
		raven.CaptureError(err, map[string]string{"type": "config"})
		return err
	}

	return nil
}

// WatchConfig 监控配置变化
func WatchConfig() error {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})
	return nil
}

// DefaultConfig 默认配置
func DefaultConfig() error {
	// basic default values
	viper.SetDefault("basic.debug", true)
	viper.SetDefault("basic.port", "8080")

	return nil
}
