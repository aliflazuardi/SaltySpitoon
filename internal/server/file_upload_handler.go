package server

import (
	"SaltySpitoon/internal/constants"
	"fmt"
	"log"
	"net/http"
)

func (s *Server) fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "POST" {
		sendErrorResponse(w, http.StatusBadRequest, "Method not allowed")
		return
	}

	if err := r.ParseMultipartForm(100 << 10); err != nil { // 100 KB
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error parsing multipart form: %v", err))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}
	defer file.Close()

	uri, err := s.service.UploadFile(ctx, file, header.Filename, header.Size)
	if err != nil {
		switch err {
		case constants.ErrMaximumFileSize:
			log.Println("invalid file size")
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		case constants.ErrInvalidFileType:
			log.Println("invalid file type")
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		default:
			log.Printf("internal server error: %v", err)
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := FileUploadResponse{
		Uri: uri,
	}

	sendResponse(w, http.StatusOK, resp)
	defer r.Body.Close()

}
