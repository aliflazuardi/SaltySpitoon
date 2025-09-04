package service

import (
	"SaltySpitoon/internal/repository"
	"context"
)

type Service struct {
	repository Repository
}

// note: not ideal, might need adapter layer because return type is defined in the repository package
type Repository interface {
	// User
	SelectUserByEmail(ctx context.Context, email string) (repository.SelectUserByEmailRow, error)
	CreateUser(ctx context.Context, arg repository.CreateUserParams) error

	// Activity
	CreateActivity(ctx context.Context, arg repository.CreateActivityParams) (repository.Activity, error)
	DeleteActivity(ctx context.Context, id int64) (int64, error)
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}
