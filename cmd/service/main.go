package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"test/internal/api"
	"test/internal/app/handler"
	"test/internal/domain/service"
	"test/internal/infrastructure/persistence/postgres/pg_repository"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	userRepo := pg_repository.NewUserRepository(db)
	teamRepo := pg_repository.NewTeamRepository(db, userRepo)
	prRepo := pg_repository.NewPrRepository(db)

	userService := service.NewUserService(userRepo)
	teamService := service.NewTeamService(teamRepo, userRepo, prRepo)
	prService := service.NewPrService(prRepo, userRepo, teamRepo)

	userHandler := handler.NewUserHandler(userService, prService)
	teamHandler := handler.NewTeamHandler(teamService)
	prHandler := handler.NewPrHandler(prService)

	apiHandler := handler.NewAPIHandler(teamHandler, userHandler, prHandler)
	r := chi.NewRouter()
	api.HandlerFromMux(apiHandler, r)
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		swagger, err := api.GetSwagger()
		if err != nil {
			http.Error(w, "Failed to load swagger spec: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(swagger)
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server running on :" + port)
	hostWithPort := fmt.Sprintf("%s:%s", hostname, port)
	err = http.ListenAndServe(hostWithPort, r)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
