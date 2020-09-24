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

	server := &grpc.Server{
		Market: &market.Market{
			UserService:    userService,
			ProductService: productService,
		},
	}

	err := server.Run(config.GRPC.Port)
	if err != nil {
		panic(err)
	}
}
