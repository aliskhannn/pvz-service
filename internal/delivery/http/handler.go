package http

import (
	"fmt"
	"github.com/aliskhannn/pvz-service/internal/config"
	"github.com/aliskhannn/pvz-service/internal/middleware"
	"github.com/aliskhannn/pvz-service/internal/usecase"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func NewRouter(
	authUC usecase.AuthUseCase,
	pvzUC usecase.PvzUseCase,
	receptionUC usecase.ReceptionUseCase,
	productUC usecase.ProductUseCase,
) http.Handler {
	r := chi.NewRouter()

	authHandler := NewAuthHandler(authUC)
	pvzHandler := NewPVZHandler(pvzUC)
	receptionHandler := NewReceptionHandler(receptionUC)
	productHandler := NewProductHandler(productUC)

	r.Post("/dummyLogin", authHandler.DummyLogin)
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)

	r.With(middleware.AuthMiddleware).Route("/pvz", func(r chi.Router) {
		r.Post("/", pvzHandler.CreatePVZ)
		r.Get("/", pvzHandler.GetAllPVZsWithReceptions)
	})

	r.With(middleware.AuthMiddleware).Route("/pvz/{pvzId}", func(r chi.Router) {
		r.Post("/close_last_reception", receptionHandler.CloseLastReception)
		r.Post("/delete_last_product", productHandler.DeleteLatProductFromReception)
	})

	r.With(middleware.AuthMiddleware).Post("/receptions", receptionHandler.CreateReception)

	r.With(middleware.AuthMiddleware).Post("/products", productHandler.AddProductToReception)

	return r
}

func Start(cfg *config.Config, r http.Handler) {
	address := fmt.Sprintf("%s", cfg.Server.HTTPPort)

	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
