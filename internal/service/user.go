package service

import (
	"SaltySpitoon/internal/model"
	"SaltySpitoon/internal/repository"
	"SaltySpitoon/internal/utils"
	"context"
	"database/sql"
	"strconv"
)

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.repository.SelectUserByEmail(ctx, email)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return "", constants.ErrUserNotFound
		}
		return "", err
	}
	if user.ID == 0 { // user not found
		return "", constants.ErrUserNotFound
	}

	if !utils.VerifyPassword(password, user.PasswordHash) {
		return "", constants.ErrUserWrongPassword
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Register(ctx context.Context, email string, password string) (string, error) {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}
	params := repository.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	}
	userID, err := s.repository.CreateUser(ctx, params)
	if err != nil {
		if utils.IsErrDBConstraint(err) {
			return "", constants.ErrEmailAlreadyExists
		}
		return "", err
	}

	token, err := utils.GenerateToken(userID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) GetProfile(ctx context.Context, id int64) (repository.SelectProfileByIdRow, error) {
	user, err := s.repository.SelectProfileById(ctx, id)
	// wrong id, error id
	if err != nil {
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
