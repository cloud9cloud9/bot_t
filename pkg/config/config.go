package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken     string
	PocketConsumerKey string
	AuthServerURI     string
	TelegramBotURI    string `mapstructure:"bot_uri"`
	DBPath            string `mapstructure:"db_file"`

	Messages Messages
}

type Messages struct {
	Responses
	Errors
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	SavedSuccessfully string `mapstructure:"saved_successfully"`
	UnknownCommand    string `mapstructure:"unknown_command"`
}

type Errors struct {
	Default      string `mapstructure:"default"`
	InvalidURI   string `mapstructure:"invalid_uri"`
	Unauthorized string `mapstructure:"unauthorized"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

func Init() (*Config, error) {
	if err := setUpViper(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := fromEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("messages.response", &cfg.Messages.Responses); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("messages.error", &cfg.Messages.Errors); err != nil {
		return err
	}

	return nil
}

func fromEnv(cfg *Config) error {

	if err := viper.BindEnv("token"); err != nil {
		return err
	}
	cfg.TelegramToken = viper.GetString("token")

	if err := viper.BindEnv("consumer_key"); err != nil {
		return err
	}
	cfg.PocketConsumerKey = viper.GetString("consumer_key")

	if err := viper.BindEnv("auth_server_uri"); err != nil {
		return err
	}
	cfg.AuthServerURI = viper.GetString("auth_server_uri")

	return nil
}

func setUpViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	return viper.ReadInConfig()
}
