package api

import "github.com/spf13/viper"

type Config struct {
	// The port to bind the web application server to
	Port int

	// The number of proxies positioned in front of the API
	ProxyCount int
}

func InitConfig() (*Config, error) {
	config := &Config{
		Port:       viper.GetInt("Port"),
		ProxyCount: viper.GetInt("ProxyCount"),
	}
	if config.Port == 0 {
		config.Port = 9095
	}
	return config, nil
}
