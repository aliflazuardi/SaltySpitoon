package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func (s *Server) deleteActivityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	prefix := "/v1/activity/"
	path := r.URL.Path
	if len(path) <= len(prefix) {
		sendErrorResponse(w, http.StatusBadRequest, "missing activityId")
		return
	}

	activityID := path[len(prefix):]
	activityIDInt, err := strconv.Atoi(activityID)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid activityId")
		return
	}

	ctx := r.Context()
	err = s.service.DeleteActivity(ctx, int64(activityIDInt))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendErrorResponse(w, http.StatusNotFound, "activityId not found")
			return
		}
		log.Println("failed to delete activity:", err)
		sendErrorResponse(w, http.StatusInternalServerError, "server error")
		return
	}

	sendResponse(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func (s *Server) patchActivityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	defer r.Body.Close()

	// Ambil activityId dari URL
	const prefix = "/v1/activity/"
	activityID := strings.TrimPrefix(r.URL.Path, prefix)
	if activityID == "" {
		sendErrorResponse(w, http.StatusBadRequest, "missing activityId")
		return
	}
	activityIDInt, err := strconv.Atoi(activityID)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid activityId")
		return
	}

	// Decode request body
	var req PatchActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid json")
		return
	}

	// Call service
	ctx := r.Context()
	res, err := s.service.PatchActivity(ctx, int64(activityIDInt), req)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			sendErrorResponse(w, http.StatusNotFound, "activityId not found")
		default:
			log.Printf("failed to patch activity (id=%d): %v", activityIDInt, err)
			sendErrorResponse(w, http.StatusInternalServerError, "server error")
		}
		return
	}

	sendResponse(w, http.StatusOK, res)
}
