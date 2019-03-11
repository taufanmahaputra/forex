package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"sync"
)

type Config struct {
	App      App
	Database Database
	Redis    Redis
}

type App struct {
	Appname string
	Port    string
}

type Database struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSL      string
}

type Redis struct {
	Host     string
	Port     string
	Password string
}

var cfg Config
var once sync.Once

func GetConfig() Config {
	once.Do(func() {
		var data string
		switch os.Getenv("ENV") {
		case "PRODUCTION":
			data = "./config/forex.production.toml"
		default:
			data = "./config/forex.development.toml"
		}

		if _, err := toml.DecodeFile(data, &cfg); err != nil {
			log.Println(err)
		}
	})

	return cfg
}
