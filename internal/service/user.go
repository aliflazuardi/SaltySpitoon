package service

import (
	"SaltySpitoon/internal/constants"
	"SaltySpitoon/internal/model"
	"SaltySpitoon/internal/repository"
	"context"
	"database/sql"
	"strconv"
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

func PtrToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false} // becomes NULL in DB
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func PtrInt64ToNullString(i *int64) sql.NullString {
	if i == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: strconv.FormatInt(*i, 10),
		Valid:  true,
	}
}
func (s *Service) PatchProfile(ctx context.Context, id int64, req model.PatchUserModel) (repository.PatchProfileByIdParams, error) {
	params := repository.PatchProfileByIdParams{
		ID:         id,
		Preference: PtrToNullString(req.Preference),
		WeightUnit: PtrToNullString(req.Weightunit),
		HeightUnit: PtrToNullString(req.Heightunit),
		Weight:     PtrInt64ToNullString(req.Weight),
		Height:     PtrInt64ToNullString(req.Height),
		Name:       PtrToNullString(req.Name),
		ImageUri:   PtrToNullString(req.Imageuri),
	}

	err := s.repository.PatchProfileById(ctx, params)
	if err != nil {
		return repository.PatchProfileByIdParams{}, err
	}

	// user, err := PatchProfile
	return params, nil
}
