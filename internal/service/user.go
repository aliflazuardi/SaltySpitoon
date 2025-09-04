package service

import (
	"SaltySpitoon/internal/constants"
	"SaltySpitoon/internal/repository"
	"context"
	"strings"
)

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	return "token", nil
}

func (s *Service) Register(ctx context.Context, email string, password string) (string, error) {
	return "token", nil
}

func (s *Service) GetProfile(ctx context.Context, id int64) (repository.SelectProfileByIdRow, error) {
	user, err := s.repository.SelectProfileById(ctx, id)
	// wrong id, error id
	if err != nil {
		if strings.Contains(err.Error(), "sql: no row in result set") {
			return repository.SelectProfileByIdRow{}, constants.ErrUserNotFound
		}
		return repository.SelectProfileByIdRow{}, err
	}

	return user, nil
}
