package main

import (
	"context"
	"fmt"
	"github.com/aliskhannn/pvz-service/internal/config"
	"github.com/aliskhannn/pvz-service/internal/delivery/http"
	"github.com/aliskhannn/pvz-service/internal/repository/postgres"
	"github.com/aliskhannn/pvz-service/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbpool.Close()

	userRepo := postgres.NewUserRepository(dbpool)
	pvzRepo := postgres.NewPVZRepository(dbpool)
	receptionRepo := postgres.NewReceptionRepository(dbpool)
	productRepo := postgres.NewProductRepository(dbpool)

	authUC := usecase.NewAuthUseCase(userRepo)
	pvzUC := usecase.NewPvzUseCase(pvzRepo)
	receptionUC := usecase.NewReceptionUseCase(receptionRepo)
	productUC := usecase.NewProductUseCase(productRepo)

	router := http.NewRouter(authUC, pvzUC, receptionUC, productUC)

	log.Printf("HTTP server running on port %s", cfg.Server.HTTPPort)
	http.Start(cfg, router)
}
