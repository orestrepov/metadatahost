package app

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// A secret string used for session cookies, passwords, etc.
	SecretKey []byte

	// SSL Labs API URL to get Host data
	SSLLabsAPIURL string

	// SSL Labs API version to get Host data
	SSLLabsAPIVersion string
}

func InitConfig() (*Config, error) {
	config := &Config{
		SecretKey:         []byte(viper.GetString("SecretKey")),
		SSLLabsAPIURL:     viper.GetString("SSLLabsAPIURL"),
		SSLLabsAPIVersion: viper.GetString("SSLLabsAPIVersion"),
	}
	if len(config.SecretKey) == 0 {
		return nil, fmt.Errorf("SecretKey must be set")
	}
	return config, nil
}
