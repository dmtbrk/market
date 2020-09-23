package config

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Port           int
	JWTAlg         string
	JWTSecret      interface{}
	UserServiceURL string

	HTTP HTTP
	GRPC GRPC
}

type HTTP struct {
	Port            int
	JWT             JWT
	UserServiceAddr string
}

type GRPC struct {
	Port int
	JWT  JWT
}

type JWT struct {
	Alg    string
	Secret interface{}
}

func FromEnv() *Config {
	httpPortString := getEnvDefault("HTTP_PORT", "8080")
	httpPort, err := strconv.Atoi(httpPortString)
	if err != nil {
		panic("cannot read PORT: " + err.Error())
	}

	grpcPortString := getEnvDefault("GRPC_PORT", "8080")
	grpcPort, err := strconv.Atoi(grpcPortString)
	if err != nil {
		panic("cannot read PORT: " + err.Error())
	}

	jwtAlg := getEnvDefault("JWT_ALG", "HS256")
	jwtSecret, err := getKey(os.Getenv("KEY_SERVICE_URL"))
	if err != nil {
		panic(fmt.Errorf("cannot get JWT secret: %w", err))
	}

	usURL := os.Getenv("USER_SERVICE_URL")

	return &Config{
		Port:           httpPort,
		JWTAlg:         jwtAlg,
		JWTSecret:      jwtSecret,
		UserServiceURL: usURL,

		HTTP: HTTP{
			Port: httpPort,
			JWT: JWT{
				Alg:    jwtAlg,
				Secret: jwtSecret,
			},
			UserServiceAddr: usURL,
		},

		GRPC: GRPC{
			Port: grpcPort,
			JWT: JWT{
				Alg:    jwtAlg,
				Secret: jwtSecret,
			},
		},
	}
}

func getEnvDefault(key string, d string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		val = d
	}
	return val
}

func getKey(url string) (*rsa.PublicKey, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("something went wrong")
	}

	key := &rsa.PublicKey{}
	err = json.NewDecoder(resp.Body).Decode(key)
	defer resp.Body.Close()

	return key, err
}
