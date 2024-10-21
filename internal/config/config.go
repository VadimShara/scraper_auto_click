package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Storage StorageConfig `yaml:"storage"`
}

type StorageConfig struct {
	Host 		string	`json:"host"`
	Port		string	`json:"port"`
	Database	string	`json:"database"`
	Auth_db		string	`json:"auth_db"`
	Username	string	`json:"username"`
	Password	string	`json:"password"`
	Collection	string	`json:"collection"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("C:/Users/vadim/scraper-first/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			fmt.Println(help)
		}
	})
	return instance
}