package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type App struct {
	PORT string `mapstructure:"PORT"`
}

type Prometheus struct {
	PROM_METRICS_FREQUENCY time.Duration `mapstructure:"PROM_METRICS_FREQUENCY"`
}

type Mongo struct {
	DB_NAME         string `mapstructure:"DB_NAME"`
	COLLECTION_NAME string `mapstructure:"COLLECTION_NAME"`
	URL             string `mapstructure:"URL"`
	CONFIG_STR      string `mapstructure:"CONFIG_STR"`
}

type Kafka struct {
	BROKERS_ADDR string `mapstructure:"BROKERS_ADDR"`
	TOPIC        string `mapstructure:"TOPIC"`
	GROUP_ID     string `mapstructure:"GROUP_ID"`
}

type WhatsappConsumer struct {
	GET_MESSAGE_LIMIT                       int           `mapstructure:"GET_MESSAGE_LIMIT"`
	GET_MESSAGE_FREQUENCY_SECONDS           time.Duration `mapstructure:"GET_MESSAGE_FREQUENCY_SECONDS"`
	UPDATE_STALE_MESSAGES_FREQUENCY_SECONDS time.Duration `mapstructure:"UPDATE_STALE_MESSAGES_FREQUENCY_SECONDS"`
	MAX_WHATSAPP_PUSH_RETRY_AVAILABLE       int           `mapstructure:"MAX_WHATSAPP_PUSH_RETRY_AVAILABLE"`
}

type WhatsappApi struct {
	TIMEOUT_SECONDS     time.Duration `mapstructure:"TIMEOUT_SECONDS"`
	RETRY_COUNT         int           `mapstructure:"RETRY_COUNT"`
	RETRY_WAIT_TIME     time.Duration `mapstructure:"RETRY_WAIT_TIME"`
	RETRY_MAX_WAIT_TIME time.Duration `mapstructure:"RETRY_MAX_WAIT_TIME"`
	BASE_URL            string        `mapstructure:"BASE_URL"`
	POST_ENDPOINT       string        `mapstructure:"POST_ENDPOINT"`
}

type OpenAi struct {
	TOKEN    string `mapstructure:"TOKEN"`
	TRAINING string `mapstructure:"TRAINING"`
}

type WhatsappAccountsConfig struct {
	NAME          string   `mapstructure:"NAME"`
	PHONE_ID      string   `mapstructure:"PHONE_ID"`
	NUMBER        string   `mapstructure:"NUMBER"`
	WA_HEADER     string   `mapstructure:"WA_HEADER"`
	IDENTIFIER    string   `mapstructure:"IDENTIFIER"`
	ORGANIZATION  string   `mapstructure:"ORGANIZATION"`
	ACCESS_TOKEN  string   `mapstructure:"ACCESS_TOKEN"`
	OUTGOING_ONLY bool     `mapstructure:"OUTGOING_ONLY"`
	CLIENTS       []string `mapstructure:"CLIENTS"`
}

type ClientsConfig struct {
	NAME         string   `mapstructure:"NAME"`
	API_KEY      string   `mapstructure:"API_KEY"`
	IDENTIFIER   string   `mapstructure:"IDENTIFIER"`
	ORGANIZATION string   `mapstructure:"ORGANIZATION"`
	INTENTS      []string `mapstructure:"INTENTS"`
}

type WaVerifyToken struct {
	TOKEN string `mapstructure:"TOKEN"`
}

type Config struct {
	APP                      App                      `mapstructure:"APP"`
	PROMETHEUS               Prometheus               `mapstructure:"PROMETHEUS"`
	MONGO                    Mongo                    `mapstructure:"MONGO"`
	KAFKA                    Kafka                    `mapstructure:"KAFKA"`
	API_MESSAGE_LIMIT        int                      `mapstructure:"API_MESSAGE_LIMIT"`
	WHATSAPP_CONSUMER        WhatsappConsumer         `mapstructure:"WHATSAPP_CONSUMER"`
	WHATSAPP_API             WhatsappApi              `mapstructure:"WHATSAPP_API"`
	WEBHOOK_FALLBACK         bool                     `mapstructure:"WEBHOOK_FALLBACK"`
	OPEN_AI                  OpenAi                   `mapstructure:"OPEN_AI"`
	CLIENTS                  []ClientsConfig          `mapstructure:"CLIENTS"`
	WHATSAPP_ACCOUNTS_CONFIG []WhatsappAccountsConfig `mapstructure:"WHATSAPP_ACCOUNTS_CONFIG"`
	WHATSAPP_VERIFY_TOKENS   []WaVerifyToken          `mapstructure:"WHATSAPP_VERIFY_TOKENS"`
}

var vp *viper.Viper

var config Config

func ConfigInit() error {
	vp = viper.New()

	var curEnv string
	env := os.Getenv("GO_ENV")
	if env == "prod" {
		curEnv = env
	} else {
		curEnv = "build"
	}
	envFileName := curEnv + ".yaml"

	workingdir, err := os.Getwd()
	if err != nil {
		return err
	}

	configFilePath := workingdir + "/config/" + envFileName
	vp.SetConfigFile(configFilePath)

	err = vp.ReadInConfig()
	if err != nil {
		return err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return err
	}

	return nil
}

func GetConfig() Config {
	return config
}

func GetViper() *viper.Viper {
	return vp
}

func GetClientInfo(clientIdentifier string) *ClientsConfig {
	for _, val := range config.CLIENTS {
		if val.IDENTIFIER == clientIdentifier {
			return &val
		}
	}
	return nil
}

func GetWhatsappAccountInfo(whatsappAccountIdentifier string) *WhatsappAccountsConfig {
	for _, val := range config.WHATSAPP_ACCOUNTS_CONFIG {
		if val.IDENTIFIER == whatsappAccountIdentifier {
			return &val
		}
	}
	return nil
}

func GetWhatsappAccountInfoUsingPhoneId(phoneId string) *WhatsappAccountsConfig {
	for _, val := range config.WHATSAPP_ACCOUNTS_CONFIG {
		if val.PHONE_ID == phoneId {
			return &val
		}
	}
	return nil
}

func GetWaTokens() []WaVerifyToken {
	return config.WHATSAPP_VERIFY_TOKENS
}
