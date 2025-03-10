package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WeatherAPIKey string
}

var config *Config

func GetConfig() *Config {
	if config == nil {

		if err := godotenv.Load(); err != nil {
			log.Println("Aviso: Não foi possível carregar o arquivo .env")
		}

		config = &Config{
			WeatherAPIKey: os.Getenv("WEATHER_API_KEY"),
		}

		if config.WeatherAPIKey == "" {
			log.Fatal("Chave da API do WeatherAPI não configurada")
		}

		log.Println("Configurações carregadas com sucesso")
	}
	return config
}
