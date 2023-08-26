package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/AthanatiusC/url-shortener/cmd/db"
	"github.com/AthanatiusC/url-shortener/config"
	"github.com/AthanatiusC/url-shortener/controller/router"
	"github.com/AthanatiusC/url-shortener/helper/logger"
	"github.com/AthanatiusC/url-shortener/internal/service"
	"github.com/AthanatiusC/url-shortener/repository"
	"github.com/redis/go-redis/v9"
)

func handleShutdown(ctx context.Context, redis *redis.Client) chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s
		logger.Info(ctx, "Server shutdown initiated")

		if err := redis.Close(); err != nil {
			logger.Error(err)
		}
		logger.Info(ctx, "Redis connection closed")

		defer func() {
			logger.Info(ctx, "Server shutdown complete")
			os.Exit(0)
		}()
		close(wait)
	}()
	return wait
}

func main() {
	ctx := context.TODO()
	defer func() {
		if err := recover(); err != nil {
			logger.Panic(err, string(debug.Stack()))
		}
	}()

	// Prerequisites
	config := config.InitConfig()
	redis := db.InitRedis(ctx, &config.Redis)

	// Repository level
	shrtnRepository := repository.NewRedisRepository(ctx, redis)

	// Service level
	shrtnService := service.NewShortenerService(ctx, config, shrtnRepository)

	// Router
	router.InitRouter(ctx, config, shrtnService)

	wait := handleShutdown(ctx, redis)
	logger.Info(fmt.Sprintf("Server started at port %d", config.Application.Port))
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Application.Host, config.Application.Port), nil)
	if err != nil {
		panic(err)
	}
	<-wait
}
