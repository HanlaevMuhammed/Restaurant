package config

import (
	"log"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port int
}

type AppConfig struct {
	Database DatabaseConfig
	Server   ServerConfig
}

var Config AppConfig

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("internal/config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Ошибка чтения config.yaml: %v", err)
	}

	Config.Database.Host = viper.GetString("database.host")
	Config.Database.Port = viper.GetInt("database.port")
	Config.Database.User = viper.GetString("database.user")
	Config.Database.Password = viper.GetString("database.password")
	Config.Database.Name = viper.GetString("database.name")
	Config.Database.SSLMode = viper.GetString("database.sslmode")

	Config.Server.Port = viper.GetInt("server.port")

}
