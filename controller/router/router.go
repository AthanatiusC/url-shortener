package router

import (
	"context"
	"net/http"

	"github.com/AthanatiusC/url-shortener/config"
	"github.com/AthanatiusC/url-shortener/controller/controller"
	"github.com/AthanatiusC/url-shortener/internal/service"
)

func InitRouter(ctx context.Context, config config.Config, shrtnService *service.ShortenerService) {
	ctx = context.WithValue(ctx, "source", "router")
	shrtnController := controller.NewShortenerController(ctx, config, shrtnService)

	http.HandleFunc("/", Handle(shrtnController.HandleShortenURL, http.MethodGet))
	http.HandleFunc("/delete", Handle(shrtnController.DeleteURL, http.MethodDelete))
	http.HandleFunc("/info", Handle(shrtnController.GetURL, http.MethodGet))
	http.HandleFunc("/shorten", Handle(shrtnController.ShortenURL, http.MethodPost))
}
