package response

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Responder struct {
	log *zap.Logger
}

func NewResponder(log *zap.Logger) *Responder {
	return &Responder{log: log}
}

func (r *Responder) JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res, err := json.Marshal(data)
	if err != nil {
		r.log.Error("Failed to marshal JSON", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if _, err := w.Write(res); err != nil {
		r.log.Error("Failed to write JSON response", zap.Error(err))
	}
}

func (r *Responder) Error(w http.ResponseWriter, status int, message string) {
	r.JSON(w, status, map[string]string{
		"error": message,
	})
}

func (r *Responder) InternalError(w http.ResponseWriter) {
	r.Error(w, http.StatusInternalServerError, "internal server error")
}
