package server

import (
	"SaltySpitoon/internal/model"
	"SaltySpitoon/internal/utils"
	"database/sql"
	"encoding/json"
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
	userID, _ := utils.GetUserIDFromCtx(ctx)
	user, err := s.service.GetProfile(ctx, userID)
	if err != nil {
		log.Println("user not found")
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
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

func (s *Server) patchProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req PatchUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("invalid patch request: %s\n", err.Error())
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = s.validator.Struct(req)
	if err != nil {
		log.Printf("invalid patch request: %s\n", err.Error())
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, _ := utils.GetUserIDFromCtx(ctx)
	params := model.PatchUserModel{
		Preference: req.Preference,
		Weightunit: req.Weightunit,
		Heightunit: req.Heightunit,
		Weight:     req.Weight,
		Height:     req.Height,
		Name:       req.Name,
		Imageuri:   req.Imageuri,
	}
	update, err := s.service.PatchProfile(ctx, userID, params)

	if err != nil {
		log.Println("error patch user: %s\n", err.Error())
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resp := PatchProfileResponse{
		Preference: toString(update.Preference),
		Weightunit: toString(update.WeightUnit),
		Heightunit: toString(update.HeightUnit),
		Weight:     toInt(update.Weight),
		Height:     toInt(update.Height),
		Name:       toString(update.Name),
		Imageuri:   toString(update.ImageUri),
	}
	sendResponse(w, http.StatusOK, resp)
}
