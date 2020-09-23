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

	m := &grpc.Client{}
	err := m.Connect(fmt.Sprintf("grpc_server:%d", config.GRPC.Port))
	if err != nil {
		panic(err)
	}

	httpSrv := httpserver.Server{
		Market:    m,
		JWTAlg:    config.JWTAlg,
		JWTSecret: config.JWTSecret,
	}
	httpSrv.Run(config.Port)

	if err := m.Disconnect(); err != nil {
		log.Println(err)
	}
}
