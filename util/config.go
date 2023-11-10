package util

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver             string
	DBSource             string
	MigrationURL         string
	HTTPServerAddress    string
	GRPCServerAddress    string
	TokenKey             string
	AccesTokenDuration   int
	RefreshTokenDuration int
	Environment          string
	RedisAddress         string
}

func LoadConfig(path string) (*Config, error) {

	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	AccesstokenDuration, err := strconv.Atoi(os.Getenv("TOKEN_ACCESS_DURATION"))
	if err != nil {
		return nil, err
	}
	RefreshTokenDuration, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DURATION"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DBDriver:             os.Getenv("DB_DRIVER"),
		DBSource:             os.Getenv("DB_SOURCE"),
		HTTPServerAddress:    os.Getenv("HTTP_SERVER_ADDRESS"),
		GRPCServerAddress:    os.Getenv("GRPC_SERVER_ADDRESS"),
		TokenKey:             os.Getenv("TOKEN_KEY"),
		Environment:          os.Getenv("ENVIRONMENT"),
		MigrationURL:         os.Getenv("MIGRATION_URL"),
		RedisAddress:         os.Getenv("REDIS_ADDRESS"),
		AccesTokenDuration:   AccesstokenDuration,
		RefreshTokenDuration: RefreshTokenDuration,
	}, nil

}
