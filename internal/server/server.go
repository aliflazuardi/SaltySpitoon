package server

import (
	"SaltySpitoon/internal/constants"
	"SaltySpitoon/internal/repository"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"

	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string) (string, error)
	CreateActivity(ctx context.Context, userID int64, req CreateActivityRequest) (repository.Activity, error)
	DeleteActivity(ctx context.Context, id int64) error
	PatchActivity(ctx context.Context, id int64, req PatchActivityRequest) (PatchActivityResponse, error)
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

	// Custom Validator For Activity Type
	NewServer.validator.RegisterValidation("activity_type_enum", func(fl validator.FieldLevel) bool {
		activityType := fl.Field().String()
		_, ok := constants.ActivityTypes[activityType]
		return ok
	})

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
