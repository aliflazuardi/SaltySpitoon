package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

func sendResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if body == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, error string) {
	resp := ErrorResponse{
		Error: error,
	}

	sendResponse(w, statusCode, resp)
}

type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type RegisterResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type CreateActivityResponse struct {
	ActivityID        int64  `json:"activityId"`
	ActivityType      string `json:"activityType"`
	DoneAt            string `json:"doneAt"`
	DurationInMinutes int32  `json:"durationInMinutes"`
	CaloriesBurned    int32  `json:"caloriesBurned"`
	CreatedAt         string `json:"createdAt,omitempty"`
	UpdatedAt         string `json:"updatedAt,omitempty"`
}

type DeleteActivityResponse struct {
	Message string `json:"message"`
}

type PatchActivityResponse struct {
	ActivityID        int64  `json:"activityId"`
	ActivityType      string `json:"activityType"`
	DoneAt            string `json:"doneAt"`
	DurationInMinutes int    `json:"durationInMinutes"`
	CaloriesBurned    int    `json:"caloriesBurned"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}
