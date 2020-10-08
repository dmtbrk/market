package main

import (
	"context"
	"fmt"
	"github.com/ortymid/market/config"
	"github.com/ortymid/market/grpc"
	"github.com/ortymid/market/http"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("getting config: %v", err)
	}

	productService := grpc.NewProductService(&grpc.UserIDAuthService{})

	grpcAddr := fmt.Sprintf("%s:%d", cfg.GRPCHost, cfg.GRPCPort)
	err = productService.Connect(context.TODO(), grpcAddr)
	if err != nil {
		log.Fatalf("Unable to connect to gRPC product service at %v: %v", grpcAddr, err)
	}

	httpServer := http.Server{
		AuthService:    http.NewJWTAuthService(cfg.JWTServiceURL),
		ProductService: productService,
	}

	httpAddr := fmt.Sprintf(":%d", cfg.HTTPPort)
	httpServer.Run(httpAddr)
}
