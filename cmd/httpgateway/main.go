package main

import (
	"fmt"
	"log"

	"github.com/ortymid/market/config"
	"github.com/ortymid/market/grpc"
	httpserver "github.com/ortymid/market/http"
)

func main() {
	config := config.FromEnv()

	marketClient := &grpc.Client{}
	err := marketClient.Connect(fmt.Sprintf("%s:%d", config.GRPC.Host, config.GRPC.Port))
	if err != nil {
		panic(err)
	}

	httpSrv := httpserver.Server{
		Market:    marketClient,
		JWTAlg:    config.JWTAlg,
		JWTSecret: config.JWTSecret,
	}
	httpSrv.Run(config.Port)

	if err := marketClient.Disconnect(); err != nil {
		log.Println(err)
	}
}
