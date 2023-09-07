package config

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	Port        string
	Private     bool
	TemplateDir string
	LogDir      string
	PasteDir    string
	Expire      int // Add an expiration setting
}

func LoadConfig() (*AppConfig, error) {
	var config AppConfig

	viper.SetConfigName("config") // name of the config file (without extension)
	viper.AddConfigPath(".")      // look in the current directory for the config file
	viper.AutomaticEnv()          // read in environment variables that match

	viper.SetDefault("Port", "8080")
	viper.SetDefault("Private", false)
	viper.SetDefault("TemplateDir", "templates")
	viper.SetDefault("LogDir", "logs")
	viper.SetDefault("PasteDir", "pastes")
	viper.SetDefault("Expire", 24) // Default expiration duration in hours

	// read the configuration file if it's present
	// goblin can also be configured through command lines during runtime
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read config file: %v\n", err)
	}

	// Unmarshal the configuration into the AppConfig struct
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
