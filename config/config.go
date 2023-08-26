package config

import (
	"encoding/json"
	"os"

	"github.com/AthanatiusC/url-shortener/helper/logger"
)

type Config struct {
	Application ApplicationConfig `json:"application"`
	Redis       RedisConfig       `json:"redis"`
}

type ApplicationConfig struct {
	Port               int    `json:"port"`
	Host               string `json:"host"`
	MaxExpiration      string `json:"max_expiration"`
	DefaultExpiration  string `json:"default_expiration"`
	GeneratedURLLength int    `json:"generated_url_length"` // Length of generated URL
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	TimeOut  int    `json:"timeout"` // Timeout in seconds
}

func InitConfig() Config {
	file, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err)
	}
	logger.Info("Config loaded")
	return config
}
