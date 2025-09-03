package service

import (
	"SaltySpitoon/internal/constants"
	"SaltySpitoon/internal/repository"
	"SaltySpitoon/internal/server"
	"context"
	"fmt"
	"time"
)

func (s *Service) CreateActivity(ctx context.Context, userID int64, req server.CreateActivityRequest) (repository.Activity, error) {
	// 1. cek activity type valid
	caloriesPerMinute, ok := constants.ActivityTypes[req.ActivityType]
	if !ok {
		return repository.Activity{}, fmt.Errorf("invalid activity type: %s", req.ActivityType)
	}

	// 2. kalkulasi calories
	totalCalories := caloriesPerMinute * req.DurationInMinutes

	// 3. parse doneAt ke time.Time
	doneAt, err := time.Parse(time.RFC3339, req.DoneAt)
	if err != nil {
		return repository.Activity{}, fmt.Errorf("invalid doneAt format: %w", err)
	}

	// 4. insert activity via repository
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

	// 5. return activity lengkap
	return activity, nil
}
