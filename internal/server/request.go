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
