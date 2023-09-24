package util

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver           string
	DBSource           string
	ServerAddress      string
	TokenKey           string
	AccesTokenDuration int
}

func LoadConfig(path string) (*Config, error) {

	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	tokenDuration, err := strconv.Atoi(os.Getenv("TOKEN_ACCESS_DURATION"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DBDriver:           os.Getenv("DB_DRIVER"),
		DBSource:           os.Getenv("DB_SOURCE"),
		ServerAddress:      os.Getenv("SERVER_ADDRESS"),
		TokenKey:           os.Getenv("TOKEN_KEY"),
		AccesTokenDuration: tokenDuration,
	}, nil

}
