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
	"strconv"
	"time"
)

const (
	NoFilterString      = ""
	NoFilterInt         = -1
	NoFilterCaloriesMax = 999999 // High number for max calories
)

var (
	NoFilterDateFrom = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	NoFilterDateTo   = time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
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

func (s *Service) PatchActivity(ctx context.Context, id int64, req server.PatchActivityRequest) (server.PatchActivityResponse, error) {
	if req.DurationInMinutes != nil && *req.DurationInMinutes < 1 {
		return server.PatchActivityResponse{}, fmt.Errorf("durationInMinutes must be >= 1")
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
		return server.PatchActivityResponse{}, fmt.Errorf("invalid doneAt format, must be ISO8601")
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
			return server.PatchActivityResponse{}, sql.ErrNoRows
		}
		return server.PatchActivityResponse{}, err
	}

	resp := server.PatchActivityResponse{
		ActivityID:        strconv.Itoa(int(row.ID)),
		ActivityType:      row.ActivityType,
		DoneAt:            row.DoneAt.Format(time.RFC3339Nano),
		DurationInMinutes: int(row.DurationMinutes),
		CaloriesBurned:    int(row.CaloriesBurned),
		CreatedAt:         utils.NullTimeToString(row.CreatedAt),
		UpdatedAt:         utils.NullTimeToString(row.UpdatedAt),
	}
	return resp, nil
}

func (s *Service) GetPaginatedActivity(ctx context.Context, userId int64, req server.GetPaginatedActivityRequest) ([]server.GetPaginatedActivityResponse, error) {
	params := parsePaginatedParams(req, userId)
	rows, err := s.repository.GetPaginatedActivity(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []server.GetPaginatedActivityResponse{}, sql.ErrNoRows
		}
		return []server.GetPaginatedActivityResponse{}, err
	}

	activities := make([]server.GetPaginatedActivityResponse, len(rows))
	for i, row := range rows {
		activities[i] = server.GetPaginatedActivityResponse{
			ActivityID:        fmt.Sprint(row.ID),
			ActivityType:      row.ActivityType,
			DoneAt:            utils.NullTimeToString(row.CreatedAt),
			DurationInMinutes: int(row.DurationMinutes),
			CaloriesBurned:    int(row.CaloriesBurned),
			CreatedAt:         utils.NullTimeToString(row.CreatedAt),
		}
	}

	return activities, nil
}

func parsePaginatedParams(req server.GetPaginatedActivityRequest, userId int64) repository.GetPaginatedActivityParams {
	params := repository.GetPaginatedActivityParams{
		UserID: userId,
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
	}

	if req.ActivityType != "" {
		params.Column2 = req.ActivityType
	} else {
		params.Column2 = NoFilterString
	}

	if req.DoneAtFrom != nil {
		params.Column3 = *req.DoneAtFrom
	} else {
		params.Column3 = NoFilterDateFrom
	}

	if req.DoneAtTo != nil {
		params.Column4 = *req.DoneAtTo
	} else {
		params.Column4 = NoFilterDateTo
	}

	if req.CaloriesBurnedMin != nil {
		params.Column5 = int32(*req.CaloriesBurnedMin)
	} else {
		params.Column5 = NoFilterInt
	}

	if req.CaloriesBurnedMax != nil {
		params.Column6 = int32(*req.CaloriesBurnedMax)
	} else {
		params.Column6 = NoFilterCaloriesMax
	}

	return params
}
