package utils

import (
	"log"

	"github.com/spf13/viper"
)

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() {
	viper.SetConfigName("config")               // name of config file (without extension)
	viper.SetConfigType("env")                  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/mattermost-bot/") // path to look for the config file in
	viper.AddConfigPath(".")                    // optionally look for config in the working directory
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			log.Fatal("Fatal error config file: %w", err)
		}
	}
}

func GetConfigValue(key string) string {
	LoadConfig()
	return viper.GetString(key)
}

func GetConfigBoolValue(key string) bool {
	LoadConfig()
	return viper.GetBool(key)
}
