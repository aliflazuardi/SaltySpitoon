package service

import (
	"SaltySpitoon/internal/constants"
	"SaltySpitoon/internal/repository"
	"SaltySpitoon/internal/server"
	"SaltySpitoon/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (s *Service) CreateActivity(ctx context.Context, userID int64, req server.CreateActivityRequest) (repository.Activity, error) {
	caloriesPerMinute, ok := constants.ActivityTypes[req.ActivityType]
	if !ok {
		return repository.Activity{}, fmt.Errorf("invalid activity type: %s", req.ActivityType)
	}

	totalCalories := caloriesPerMinute * req.DurationInMinutes

	doneAt, err := time.Parse(time.RFC3339, req.DoneAt)
	if err != nil {
		return repository.Activity{}, fmt.Errorf("invalid doneAt format: %w", err)
	}

	activity, err := s.repository.CreateActivity(ctx, repository.CreateActivityParams{
		UserID:          userID,
		ActivityType:    req.ActivityType,
		DoneAt:          doneAt,
		DurationMinutes: int32(req.DurationInMinutes),
		CaloriesBurned:  int32(totalCalories),
	})
	if err != nil {
		return repository.Activity{}, err
	}

	return activity, nil
}

func (s *Service) DeleteActivity(ctx context.Context, id int64) error {
	rows, err := s.repository.DeleteActivity(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Service) PatchActivity(ctx context.Context, id int64, req server.PatchActivityRequest) (repository.PatchActivityRow, error) {
	if req.DurationInMinutes != nil && *req.DurationInMinutes < 1 {
		return repository.PatchActivityRow{}, fmt.Errorf("durationInMinutes must be >= 1")
	}

	var calories *int
	if req.DurationInMinutes != nil && req.ActivityType != nil {
		if met, ok := constants.ActivityTypes[*req.ActivityType]; ok {
			c := met * (*req.DurationInMinutes)
			calories = &c
		}
	}

	doneAt, err := utils.ToNullTimeFromString(req.DoneAt)
	if err != nil {
		return repository.PatchActivityRow{}, fmt.Errorf("invalid doneAt format, must be ISO8601")
	}

	row, err := s.repository.PatchActivity(ctx, repository.PatchActivityParams{
		ID:              id,
		ActivityType:    utils.ToNullString(req.ActivityType),
		DoneAt:          doneAt,
		DurationMinutes: utils.ToNullInt32(req.DurationInMinutes),
		CaloriesBurned:  utils.ToNullInt32(calories),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.PatchActivityRow{}, sql.ErrNoRows
		}
		return repository.PatchActivityRow{}, err
	}

	return row, nil
}
