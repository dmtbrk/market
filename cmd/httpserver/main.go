package main

import (
	"database/sql"
	"fmt"
	"github.com/ortymid/market/config"
	"github.com/ortymid/market/http"
	"github.com/ortymid/market/market/product"
	"github.com/ortymid/market/storage/postgres"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(fmt.Errorf("getting config: %w", err))
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to postgres: %v", err)
	}
	log.Printf("Connected to postgres at %v", cfg.DatabaseURL)

	productService := &product.Service{
		Storage: postgres.NewProductStorage(db, "products"),
	}

	httpServer := http.Server{
		AuthService:    http.NewJWTAuthService(cfg.JWTServiceURL),
		ProductService: productService,
	}

	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	httpServer.Run(addr)
}
