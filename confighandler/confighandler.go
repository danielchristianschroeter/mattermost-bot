package confighandler

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var version = "DEV"
var App = &Application{}

var Settings Config

// application struct to hold the dependencies for our bot
type Application struct {
	Config                    Config
	Logger                    zerolog.Logger
	MattermostClient          *model.Client4
	MattermostWebSocketClient *model.WebSocketClient
	MattermostUser            *model.User
	MattermostChannel         *model.Channel
	MattermostTeam            *model.Team
}

type Config struct {
	MB_ROLES                        string
	MB_MATTERMOST_URL               string
	MB_MATTERMOST_USER              string
	MB_MATTERMOST_USERTOKEN         string
	MB_MATTERMOST_TEAM              string
	MB_MATTERMOST_CHANNEL           string
	MB_PRIVATEBIN_ENABLE            string
	MB_PRIVATEBIN_URL               string
	MB_PRIVATEBIN_FORMATTER         string
	MB_PRIVATEBIN_EXPIRE            string
	MB_PRIVATEBIN_OPENDISCUSSION    string
	MB_PRIVATEBIN_BURNAFTERREADING  string
	MB_PRIVATEBIN_PASSWORD          string
	MB_TERRAFORM_EXECUTABLE         string
	MB_TERRAFORM_WRAPPER_EXECUTABLE string
	MB_OC_EXECUTABLE                string
	MB_GOVC_EXECUTABLE              string
	MB_GOVC_URL                     string
	MB_GOVC_DATACENTER              string
	MB_GOVC_USERNAME                string
	MB_GOVC_PASSWORD                string
}

func Init() {
	App = &Application{
		Logger: zerolog.New(
			os.Stdout,
		).With().Timestamp().Logger(),
	}

	App.Config = LoadConfig()
}

func LoadConfig() Config {
	App.Logger.Info().Msg("Mattermost-Bot (v" + version + ") started.")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("env")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			log.Fatal("Fatal error config file: %w", err)
		}
	}
	// If MB_DEBUG is true, enable debug mode
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if GetConfigBoolValue("MB_DEBUG") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	settings := map[string]string{
		"MB_ROLES":                        GetConfigValue("MB_ROLES"),
		"MB_MATTERMOST_URL":               GetConfigValue("MB_MATTERMOST_URL"),
		"MB_MATTERMOST_USERTOKEN":         GetConfigValue("MB_MATTERMOST_USERTOKEN"),
		"MB_MATTERMOST_TEAM":              GetConfigValue("MB_MATTERMOST_TEAM"),
		"MB_MATTERMOST_CHANNEL":           GetConfigValue("MB_MATTERMOST_CHANNEL"),
		"MB_PRIVATEBIN_ENABLE":            GetConfigValue("MB_PRIVATEBIN_ENABLE"),
		"MB_PRIVATEBIN_URL":               GetConfigValue("MB_PRIVATEBIN_URL"),
		"MB_PRIVATEBIN_FORMATTER":         GetConfigValue("MB_PRIVATEBIN_FORMATTER"),
		"MB_PRIVATEBIN_EXPIRE":            GetConfigValue("MB_PRIVATEBIN_EXPIRE"),
		"MB_PRIVATEBIN_OPENDISCUSSION":    GetConfigValue("MB_PRIVATEBIN_OPENDISCUSSION"),
		"MB_PRIVATEBIN_BURNAFTERREADING":  GetConfigValue("MB_PRIVATEBIN_BURNAFTERREADING"),
		"MB_PRIVATEBIN_PASSWORD":          GetConfigValue("MB_PRIVATEBIN_PASSWORD"),
		"MB_TERRAFORM_EXECUTABLE":         GetConfigValue("MB_TERRAFORM_EXECUTABLE"),
		"MB_TERRAFORM_WRAPPER_EXECUTABLE": GetConfigValue("MB_TERRAFORM_WRAPPER_EXECUTABLE"),
		"MB_OC_EXECUTABLE":                GetConfigValue("MB_OC_EXECUTABLE"),
		"MB_GOVC_EXECUTABLE":              GetConfigValue("MB_GOVC_EXECUTABLE"),
		"MB_GOVC_URL":                     GetConfigValue("MB_GOVC_URL"),
		"MB_GOVC_DATACENTER":              GetConfigValue("MB_GOVC_DATACENTER"),
		"MB_GOVC_USERNAME":                GetConfigValue("MB_GOVC_USERNAME"),
		"MB_GOVC_PASSWORD":                GetConfigValue("MB_GOVC_PASSWORD"),
	}

	// Output all keys and values of the current loaded config
	message := "Current configuration: "
	for key, value := range settings {
		if key == "MB_MATTERMOST_USERTOKEN" || key == "MB_PRIVATEBIN_PASSWORD" || key == "MB_GOVC_PASSWORD" {
			continue
		}
		message += key + "=" + value + " "
	}
	message = message[:len(message)-1]

	App.Logger.Info().Msg(message)

	Settings := Config{
		MB_ROLES:                        settings["MB_ROLES"],
		MB_MATTERMOST_URL:               settings["MB_MATTERMOST_URL"],
		MB_MATTERMOST_USERTOKEN:         settings["MB_MATTERMOST_USERTOKEN"],
		MB_MATTERMOST_TEAM:              settings["MB_MATTERMOST_TEAM"],
		MB_MATTERMOST_CHANNEL:           settings["MB_MATTERMOST_CHANNEL"],
		MB_PRIVATEBIN_ENABLE:            settings["MB_PRIVATEBIN_ENABLE"],
		MB_PRIVATEBIN_URL:               settings["MB_PRIVATEBIN_URL"],
		MB_PRIVATEBIN_FORMATTER:         settings["MB_PRIVATEBIN_FORMATTER"],
		MB_PRIVATEBIN_EXPIRE:            settings["MB_PRIVATEBIN_EXPIRE"],
		MB_PRIVATEBIN_OPENDISCUSSION:    settings["MB_PRIVATEBIN_OPENDISCUSSION"],
		MB_PRIVATEBIN_BURNAFTERREADING:  settings["MB_PRIVATEBIN_BURNAFTERREADING"],
		MB_PRIVATEBIN_PASSWORD:          settings["MB_PRIVATEBIN_PASSWORD"],
		MB_TERRAFORM_EXECUTABLE:         settings["MB_TERRAFORM_EXECUTABLE"],
		MB_TERRAFORM_WRAPPER_EXECUTABLE: settings["MB_TERRAFORM_WRAPPER_EXECUTABLE"],
		MB_OC_EXECUTABLE:                settings["MB_OC_EXECUTABLE"],
		MB_GOVC_EXECUTABLE:              settings["MB_GOVC_EXECUTABLE"],
		MB_GOVC_URL:                     settings["MB_GOVC_URL"],
		MB_GOVC_DATACENTER:              settings["MB_GOVC_DATACENTER"],
		MB_GOVC_USERNAME:                settings["MB_GOVC_USERNAME"],
		MB_GOVC_PASSWORD:                settings["MB_GOVC_PASSWORD"],
	}

	return Settings
}

func checkParams(params ...string) bool {
	for _, param := range params {
		if !viper.IsSet(param) {
			return false
		}
	}
	return true
}

func GetConfigValue(key string) string {
	//LoadConfig()
	return viper.GetString(key)
}

func GetConfigBoolValue(key string) bool {
	//LoadConfig()
	return viper.GetBool(key)
}
