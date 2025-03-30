package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status string      `json:"status"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK(data interface{}) Response {
	return Response{
		Status: StatusOK,
		Data:   data,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
