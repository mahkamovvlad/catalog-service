package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/mahkamovvlad/catalog-service/internal/app/config/section"
)

type Config struct {
	Repository section.Repository
	Monitor    section.Monitor
	Processor  section.Processor
}

var Root Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Printf("инфо: .env файл не найден, используются системные переменные: %v", err)
	}

	err := envconfig.Process("APP", &Root)
	if err != nil {
		log.Fatalf("критическая ошибка конфигурации: %v", err)
	}
}
