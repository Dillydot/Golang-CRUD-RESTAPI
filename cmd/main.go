package main

import (
	"context"
	"log"
	"net/http"

	"github.com/simpler-tha/internal/app"
	"github.com/simpler-tha/internal/config"
	infrahttp "github.com/simpler-tha/internal/infra/http"
	"github.com/simpler-tha/internal/infra/postgresql"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	client, err := postgresql.NewClient(ctx, cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to initialize postgresql client: %v", err)
	}
	defer func() {
		if err := client.Conn.Close(ctx); err != nil {
			log.Printf("failed to close postgresql connection: %v", err)
		}
	}()

	productsRepository, err := postgresql.NewRepository(client)
	if err != nil {
		log.Fatalf("failed to initialize products repository: %v", err)
	}

	service, err := app.NewService(productsRepository)
	if err != nil {
		log.Fatalf("failed to initialize service: %v", err)
	}

	router, err := infrahttp.NewRouter(service)
	if err != nil {
		log.Fatalf("failed to initialize HTTP router: %v", err)
	}
	router.RegisterRoutes()

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("server error: %s", err)
	}
}
