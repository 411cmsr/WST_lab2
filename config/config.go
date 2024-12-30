package config

import (
	"WST_lab1_server_new1/internal/models"
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Структура конфигурации
type Config struct {
	GeneralServer GeneralServerConfig `yaml:"generalServer"`
	HTTPServer    HTTPServerConfig    `yaml:"httpServer"`
	Database      DatabaseConfig      `yaml:"database"`
}

// Структура конфигурации сервера
type GeneralServerConfig struct {
	Env      string          `yaml:"env" env-required:"true"`
	LogLevel string          `yaml:"logLevel" env-default:"debug"`
	DataSet  []models.Person `yaml:"persons"`
}

// Структура конфигурации HTTP сервера
type HTTPServerConfig struct {
	RunMode      string        `yaml:"runMode"`
	BindAddr     string        `yaml:"bindAddr"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

// Структура конфигурации подключения к базе данных
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	SSLMode  string `yaml:"sslMode"`
}

// Переменные конфигурации
var (
	config               Config
	GeneralServerSetting = &GeneralServerConfig{}
	HTTPServerSetting    = &HTTPServerConfig{}
	DatabaseSetting      = &DatabaseConfig{}
)

// Функция инициализации конфигурации
func Init() {
	var pathConfigFile string
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}
	//Проверяем hostname для загрузки нужной конфигурации
	if hostname == "test-XWPC" {
		pathConfigFile = "config/vm.yaml"
	} else {
		pathConfigFile = "config/pc.yaml"
	}
	//Открываем файл конфигурации
	file, err := os.Open(pathConfigFile)
	if err != nil {
		log.Fatal("error opening file config", zap.Error(err))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("error closing file config", zap.Error(err))
		}
	}(file)
	//Читаем файл конфигурации
	decoder := yaml.NewDecoder(file)
	//Привязываем переменные конфигурации
	if err := decoder.Decode(&config); err != nil {
		log.Fatal("error decoding file config", zap.Error(err))
	}
	*GeneralServerSetting = config.GeneralServer
	*HTTPServerSetting = config.HTTPServer
	*DatabaseSetting = config.Database

}
