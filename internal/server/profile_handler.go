package server

import (
	"SaltySpitoon/internal/model"
	"SaltySpitoon/internal/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func rawJSONHasField(body io.ReadCloser, field string) bool {
	defer body.Close()
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(body).Decode(&raw); err != nil {
		return false
	}
	_, ok := raw[field]
	return ok
}

func (s *Server) getProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := utils.GetUserIDFromCtx(ctx)
	user, err := s.service.GetProfile(ctx, userID)
	if err != nil {
		log.Println("user not found", err.Error())
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resp := GetProfileResponse{
		Preference: toString(user.Preference),
		Weightunit: toString(user.Weightunit),
		Heightunit: toString(user.Heightunit),
		Weight:     int32toInt(user.Weight),
		Height:     int32toInt(user.Height),
		Email:      user.Email,
		Name:       toString(user.Name),
		Imageuri:   toString(user.Imageuri),
	}
	sendResponse(w, http.StatusOK, resp)
}

func (s *Server) patchProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contentTypeHeader := r.Header.Get("Content-Type")
	if !strings.EqualFold(contentTypeHeader, "application/json") {
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	var req PatchUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("invalid payload: %s\n", err.Error())
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = s.validator.Struct(req)
	if err != nil {
		log.Printf("invalid vaidator: %s\n", err.Error())
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if rawJSONHasField(r.Body, "imageUri") && req.Imageuri == nil {
		log.Printf("invalid: name explicitly null")
		sendErrorResponse(w, http.StatusBadRequest, "name must not be null")
		return
	}

	if req.Imageuri != nil {
		strUri := *req.Imageuri
		log.Printf("struri %s\n", strUri)
		if strUri == "" {
			log.Printf("invalid validator")
			sendErrorResponse(w, http.StatusBadRequest, "bad request")
			return
		}
		if !strings.Contains(strUri, ".com") &&
			!strings.Contains(strUri, ".org") &&
			!strings.Contains(strUri, ".net") &&
			!strings.Contains(strUri, ".io") &&
			!strings.Contains(strUri, ".co") &&
			!strings.Contains(strUri, ".xyz") {
			log.Printf("invalid validator")
			sendErrorResponse(w, http.StatusBadRequest, "bad request")
			return
		}
	}

	// Catch case where name is explicitly set to null (i.e., req.Name == nil, but field present)
	if rawJSONHasField(r.Body, "name") && req.Name == nil {
		log.Printf("invalid: name explicitly null")
		sendErrorResponse(w, http.StatusBadRequest, "name must not be null")
		return
	}

	if req.Name != nil {
		strName := *req.Name
		log.Printf("strname %s\n", strName)
		if strName == "" {
			log.Printf("invalid validator name")
			sendErrorResponse(w, http.StatusBadRequest, "bad request")
			return
		}
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
		log.Printf("error patch user: %s\n", err.Error())
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resp := PatchProfileResponse{
		Preference: toString(update.Preference),
		Weightunit: toString(update.WeightUnit),
		Heightunit: toString(update.HeightUnit),
		Weight:     int32toInt(update.Weight),
		Height:     int32toInt(update.Height),
		Name:       toString(update.Name),
		Imageuri:   toString(update.ImageUri),
	}
	sendResponse(w, http.StatusOK, resp)
}
