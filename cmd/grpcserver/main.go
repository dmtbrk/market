package main

import (
	"github.com/ortymid/market/config"
	"github.com/ortymid/market/grpc"
	"github.com/ortymid/market/market"
	"github.com/ortymid/market/service/http"
	"github.com/ortymid/market/service/mem"
)

func main() {
	config := config.FromEnv()

	userService := http.NewUserService(config.UserServiceURL)
	productService := mem.NewProductService()

	m := &market.Market{
		UserService:    userService,
		ProductService: productService,
	}

	server := &grpc.Server{
		Market: m,

		JWTAlg:    config.GRPC.JWT.Alg,
		JWTSecret: config.GRPC.JWT.Secret,
	}

	err := server.Run(config.GRPC.Port)
	if err != nil {
		panic(err)
	}
}
