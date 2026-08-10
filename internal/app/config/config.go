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
	_ = godotenv.Load()
	err := envconfig.Process("APP", &Root)
	if err != nil {
		log.Fatal("ошибка при загрузке конфигурации: ", err)
	}
}
