package response

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AthanatiusC/url-shortener/helper/logger"
	"github.com/AthanatiusC/url-shortener/model"
)

func Success(ctx context.Context, w http.ResponseWriter, data interface{}, status int, message ...interface{}) {
	defer logger.InfoContext(ctx, "response success returned", data)
	var response model.ShortenerResponse
	response.Status = "success"
	response.Data = data
	for _, msg := range message {
		response.Message += msg.(string)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func Error(ctx context.Context, w http.ResponseWriter, httpCode int, message ...interface{}) {
	defer logger.InfoContext(ctx, "response error returned", message)
	if httpCode == 0 {
		httpCode = http.StatusInternalServerError
	}
	var response model.ShortenerResponse
	response.Status = "error"
	for _, msg := range message {
		response.Message += msg.(string)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(response)
}
