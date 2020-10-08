package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ortymid/market/config"
	"github.com/ortymid/market/grpc"
	"github.com/ortymid/market/market/product"
	"github.com/ortymid/market/storage/mongo"
	"github.com/ortymid/market/storage/postgres"
	"github.com/ortymid/market/storage/redis"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.FromEnv()
	if err != nil {
		return fmt.Errorf("getting config: %w", err)
	}

	productStorage, err := getProductStorage(cfg)
	if err != nil {
		return fmt.Errorf("unable to get product storage: %w", err)
	}

	grpcServer := grpc.Server{
		AuthService: &grpc.UserIDAuthService{},
		ProductService: &product.Service{
			Storage: productStorage,
		},
	}

	addr := fmt.Sprintf(":%d", cfg.GRPCPort)
	return grpcServer.Run(addr)
}

func getProductStorage(cfg *config.Config) (product.Storage, error) {
	if len(cfg.DatabaseURL) == 0 {
		return nil, errors.New("database url is empty")
	}

	db := strings.SplitN(cfg.DatabaseURL, ":", 2)[0]

	switch db {
	case "redis":
		return getRedisProductStorage(cfg.DatabaseURL)
	case "postgres":
		return getPostgresProductStorage(cfg.DatabaseURL)
	case "mongodb":
		return getMongoProductStorage(cfg.DatabaseURL)
	default:
		return nil, errors.New("unknown database in database url")
	}
}

func getRedisProductStorage(dbURL string) (*redis.ProductStorage, error) {
	db, err := redis.NewClientFromURL(dbURL)
	if err != nil {
		return nil, err
	}

	return redis.NewProductStorage(db, "products"), nil
}

func getPostgresProductStorage(dbURL string) (*postgres.ProductStorage, error) {
	db, err := postgres.NewDBFromURL(dbURL)
	if err != nil {
		return nil, err
	}

	return postgres.NewProductStorage(db, "products"), nil
}

func getMongoProductStorage(dbURL string) (*mongo.ProductStorage, error) {
	client, err := mongo.NewClientFromURL(dbURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	col := client.Database("market").Collection("products")

	return mongo.NewProductStorage(col), nil
}
