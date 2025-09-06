package server

import (
	"SaltySpitoon/internal/constants"
	"errors"
	"strings"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type PatchUserRequest struct {
	Preference *string `json:"preference,omitempty" validate:"required,oneof=CARDIO WEIGHT"`
	Weightunit *string `json:"weightUnit,omitempty" validate:"required,oneof=KG LBS"`
	Heightunit *string `json:"heightUnit,omitempty" validate:"required,oneof=CM INCH"`
	Weight     *int    `json:"weight,omitempty" validate:"required,gte=10,lte=1000"`
	Height     *int    `json:"height,omitempty" validate:"required,gte=3,lte=250"`
	Name       *string `json:"name,omitempty" validate:"min=2,max=60"`
	Imageuri   *string `json:"imageUri,omitempty" validate:"uri"`
}

type CreateActivityRequest struct {
	ActivityType      string `json:"activityType" validate:"required,activity_type_enum"`
	DoneAt            string `json:"doneAt" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	DurationInMinutes int    `json:"durationInMinutes" validate:"required,min=1"`
}

type PatchActivityRequest struct {
	ActivityType      *string `json:"activityType,omitempty"`
	DoneAt            *string `json:"doneAt,omitempty"`
	DurationInMinutes *int    `json:"durationInMinutes,omitempty"`
}

func (r *PatchActivityRequest) Validate() error {
	if r.ActivityType == nil || strings.TrimSpace(*r.ActivityType) == "" {
		return errors.New("activityType is required")
	}
	if _, ok := constants.ActivityTypes[*r.ActivityType]; !ok {
		return errors.New("invalid activityType")
	}

	if r.DoneAt == nil || strings.TrimSpace(*r.DoneAt) == "" {
		return errors.New("doneAt is required")
	}
	if _, err := time.Parse(time.RFC3339, *r.DoneAt); err != nil {
		return errors.New("invalid doneAt format")
	}

	if r.DurationInMinutes == nil || *r.DurationInMinutes < 1 {
		return errors.New("durationInMinutes must be >= 1")
	}
	return nil
}

type GetPaginatedActivityRequest struct {
	Limit             int
	Offset            int
	ActivityType      string
	DoneAtFrom        *time.Time
	DoneAtTo          *time.Time
	CaloriesBurnedMin *int
	CaloriesBurnedMax *int
}
