package server

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"os"
	"strconv"
	"time"

	"SaltySpitoon/internal/model"
	"SaltySpitoon/internal/repository"

	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string) (string, error)
	GetProfile(ctx context.Context, id int64) (repository.SelectProfileByIdRow, error)
	PatchProfile(ctx context.Context, id int64, req model.PatchUserModel) (repository.PatchProfileByIdParams, error)
}

type Server struct {
	port      int
	service   Service
	validator *validator.Validate
}

func NewServer(service Service) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:      port,
		service:   service,
		validator: validator.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
