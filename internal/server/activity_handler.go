package server

import (
	"SaltySpitoon/internal/utils"
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

	ct := r.Header.Get("Content-Type")
	if ct == "" || !strings.HasPrefix(ct, "application/json") {
		sendErrorResponse(w, http.StatusBadRequest, "invalid content type")
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

	userID, ok := utils.GetUserIDFromCtx(ctx)
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

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
		ActivityID:        strconv.Itoa(int(activity.ID)),
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

	idStr := strings.TrimPrefix(r.URL.Path, "/v1/activity/")
	if idStr == "" {
		sendErrorResponse(w, http.StatusNotFound, "activityId not found")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "activityId not found")
		return
	}

	ctx := r.Context()
	err = s.service.DeleteActivity(ctx, int64(id))
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

	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		sendErrorResponse(w, http.StatusBadRequest, "invalid content type")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/v1/activity/")
	if idStr == "" {
		sendErrorResponse(w, http.StatusNotFound, "activityId not found")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "activityId not found")
		return
	}

	var req PatchActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	row, err := s.service.PatchActivity(ctx, int64(id), req)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			sendErrorResponse(w, http.StatusNotFound, "activityId not found")
		default:
			log.Printf("failed to patch activity (id=%d): %v", id, err)
			sendErrorResponse(w, http.StatusInternalServerError, "server error")
		}
		return
	}

	resp := PatchActivityResponse{
		ActivityID:        strconv.Itoa(int(row.ID)),
		ActivityType:      row.ActivityType,
		DoneAt:            row.DoneAt.Format(time.RFC3339Nano),
		DurationInMinutes: int(row.DurationMinutes),
		CaloriesBurned:    int(row.CaloriesBurned),
		CreatedAt:         utils.NullTimeToString(row.CreatedAt),
		UpdatedAt:         utils.NullTimeToString(row.UpdatedAt),
	}
	sendResponse(w, http.StatusOK, resp)
}
