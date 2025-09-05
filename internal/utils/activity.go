package utils

import (
	"SaltySpitoon/internal/constants"
	"fmt"
	"time"
)

func IsValidActivityType(activityType string) bool {
	validTypes := constants.ActivityTypes
	_, ok := validTypes[activityType]
	return ok
}

func ParseISODate(dateStr string) (time.Time, error) {
	// Try multiple ISO formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
}
