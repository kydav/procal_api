package routes

import (
	"encoding/json"
	"net/http"

	"log/slog"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func returnSuccess(writer http.ResponseWriter, responseData interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(responseData); err != nil {
		slog.Warn("Error encoding response")
		return
	}
}

func returnError(writer http.ResponseWriter, errorMessage string, httpStatus int, err error) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpStatus)
	var resp ErrorResponse

	if err != nil {
		resp.Error = err.Error()
		resp.Message = errorMessage
	} else {
		resp.Error = errorMessage
	}
	jsonResp, _ := json.Marshal(resp)
	_, err = writer.Write(jsonResp)
	if err != nil {
		slog.Warn("Error writing response")
	}
}
