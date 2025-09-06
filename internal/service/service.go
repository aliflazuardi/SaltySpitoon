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
	CreateUser(ctx context.Context, arg repository.CreateUserParams) (int64, error)
	SelectProfileById(ctx context.Context, id int64) (repository.SelectProfileByIdRow, error)
	PatchProfileById(ctx context.Context, params repository.PatchProfileByIdParams) error

	// Activity
	CreateActivity(ctx context.Context, arg repository.CreateActivityParams) (repository.Activity, error)
	DeleteActivity(ctx context.Context, id int64) (int64, error)
	PatchActivity(ctx context.Context, arg repository.PatchActivityParams) (repository.PatchActivityRow, error)
	GetPaginatedActivity(ctx context.Context, arg repository.GetPaginatedActivityParams) ([]repository.GetPaginatedActivityRow, error)
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}
