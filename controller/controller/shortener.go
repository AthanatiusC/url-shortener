package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AthanatiusC/url-shortener/config"
	"github.com/AthanatiusC/url-shortener/helper/logger"
	"github.com/AthanatiusC/url-shortener/helper/response"
	"github.com/AthanatiusC/url-shortener/internal/service"
	"github.com/AthanatiusC/url-shortener/model"
	"github.com/google/uuid"
)

type ShortenerController struct {
	config  config.Config
	service *service.ShortenerService
}

func NewShortenerController(ctx context.Context, config config.Config, shrtnService *service.ShortenerService) *ShortenerController {
	return &ShortenerController{
		config:  config,
		service: shrtnService,
	}
}

func (c *ShortenerController) HandleShortenURL(w http.ResponseWriter, r *http.Request) {
	requestID, _ := uuid.NewUUID()
	ctx := context.WithValue(r.Context(), "request_id", requestID)
	logger.LogRequest(ctx, r, "HandleShortenURL")

	var url string
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) >= 2 {
		url = pathParts[1]
	} else {
		response.Error(ctx, w, http.StatusNotFound, fmt.Errorf("invalid URL"))
		return
	}

	result, err := c.service.HandleShortenURL(ctx, url)
	if err != nil {
		response.Error(ctx, w, result.HttpCode, err.Error())
		return
	}
	http.Redirect(w, r, result.URL, http.StatusMovedPermanently)
}

func (c *ShortenerController) GetURL(w http.ResponseWriter, r *http.Request) {
	requestID, _ := uuid.NewUUID()
	ctx := context.WithValue(r.Context(), "request_id", requestID)
	logger.LogRequest(ctx, r, "GetURL")

	url := r.URL.Query().Get("short_url")
	result, err := c.service.GetURL(ctx, url)
	if err != nil {
		response.Error(ctx, w, result.HttpCode, err.Error())
		return
	}
	response.Success(ctx, w, result, http.StatusOK)
}

func (c *ShortenerController) ShortenURL(w http.ResponseWriter, r *http.Request) {
	requestID, _ := uuid.NewUUID()
	ctx := context.WithValue(r.Context(), "request_id", requestID)
	logger.LogRequest(ctx, r, "ShortenURL")

	var request model.ShortenerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Error(ctx, w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := c.service.ShortenURL(ctx, request.URL, request.Exp)
	if err != nil {
		response.Error(ctx, w, result.HttpCode, err.Error())
		return
	}
	response.Success(ctx, w, result, http.StatusOK)
}

func (c *ShortenerController) DeleteURL(w http.ResponseWriter, r *http.Request) {
	requestID, _ := uuid.NewUUID()
	ctx := context.WithValue(r.Context(), "request_id", requestID)
	logger.LogRequest(ctx, r, "DeleteURL")

	url := r.URL.Query().Get("short_url")
	httpCode, err := c.service.DeleteURL(ctx, url)
	if err != nil {
		response.Error(ctx, w, httpCode, err.Error())
		return
	}
	response.Success(ctx, w, nil, http.StatusOK, "URL deleted successfully")
}
