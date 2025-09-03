package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (s *Server) createActivityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ctx := r.Context()
	var req CreateActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := s.validator.Struct(req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// userID, ok := utils.GetUserIDFromCtx(ctx)
	// if !ok {
	// 	sendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
	// 	return
	// }

	// sementara hardcode userID
	userID := int64(1)

	activity, err := s.service.CreateActivity(ctx, userID, req)
	if err != nil {
		log.Println("failed to create activity:", err)
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	createdAt, updatedAt := "", ""
	if activity.CreatedAt.Valid {
		createdAt = activity.CreatedAt.Time.Format(time.RFC3339)
	}
	if activity.UpdatedAt.Valid {
		updatedAt = activity.UpdatedAt.Time.Format(time.RFC3339)
	}

	resp := CreateActivityResponse{
		ActivityID:        activity.ID,
		ActivityType:      activity.ActivityType,
		DoneAt:            activity.DoneAt.Format(time.RFC3339),
		DurationInMinutes: activity.DurationMinutes,
		CaloriesBurned:    activity.CaloriesBurned,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	sendResponse(w, http.StatusCreated, resp)
}
