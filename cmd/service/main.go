package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"test/internal/api"
	"test/internal/app/handler"
	"test/internal/domain/service"
	"test/internal/infrastructure/persistence/postgres/pg_repository"
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
	server := api.Handler(apiHandler)

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
	err = http.ListenAndServe(hostWithPort, server)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
