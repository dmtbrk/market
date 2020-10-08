package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	HTTPHost string
	HTTPPort int

	GRPCHost string
	GRPCPort int

	JWTServiceURL string

	DatabaseURL string
}

func FromEnv() (*Config, error) {
	httpHost := os.Getenv("MARKET_HTTP_HOST")

	httpPortString := os.Getenv("MARKET_HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortString)
	if err != nil {
		return nil, fmt.Errorf("parsing HTTP_PORT: %w", err)
	}

	grpcHost := os.Getenv("MARKET_GRPC_HOST")

	grpcPortString := os.Getenv("MARKET_GRPC_PORT")
	grpcPort, err := strconv.Atoi(grpcPortString)
	if err != nil {
		return nil, fmt.Errorf("parsing GRPC_PORT: %w", err)
	}

	jwtServiceURL := os.Getenv("MARKET_JWT_SERVICE_URL")

	databaseURL := os.Getenv("MARKET_DATABASE_URL")

	return &Config{
		HTTPHost: httpHost,
		HTTPPort: httpPort,

		GRPCHost: grpcHost,
		GRPCPort: grpcPort,

		JWTServiceURL: jwtServiceURL,

		DatabaseURL: databaseURL,
	}, nil
}
