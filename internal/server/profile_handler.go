package server

import (
	"SaltySpitoon/internal/constants"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func toString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func toInt(ni sql.NullString) int64 {
	if ni.Valid {
		if i, err := strconv.ParseInt(ni.String, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

func (s *Server) getProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := int64(12345)
	user, err := s.service.GetProfile(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrUserNotFound) {
			sendErrorResponse(w, http.StatusNotFound, fmt.Sprintf("user not found"))
		}
		log.Println("user not found")
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	resp := GetProfileResponse{
		Preference: toString(user.Preference),
		Weightunit: toString(user.Weightunit),
		Heightunit: toString(user.Heightunit),
		Weight:     toInt(user.Weight),
		Height:     toInt(user.Height),
		Email:      user.Email,
		Name:       toString(user.Name),
		Imageuri:   toString(user.Imageuri),
	}
	sendResponse(w, http.StatusOK, resp)
}
