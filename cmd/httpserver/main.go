package main

import (
	"github.com/ortymid/market/config"
	httpserver "github.com/ortymid/market/http"
	"github.com/ortymid/market/market"
	httpservice "github.com/ortymid/market/service/http"
	"github.com/ortymid/market/service/mem"
)

func main() {
	config := config.FromEnv()

	userService := httpservice.NewUserService(config.UserServiceURL)
	productService := mem.NewProductService()
	m := &market.Market{
		UserService:    userService,
		ProductService: productService,
	}

	httpSrv := httpserver.Server{
		Market:    m,
		JWTAlg:    config.JWTAlg,
		JWTSecret: config.JWTSecret,
	}
	httpSrv.Run(config.Port)
}
