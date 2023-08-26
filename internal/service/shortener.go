package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AthanatiusC/url-shortener/config"
	"github.com/AthanatiusC/url-shortener/helper/logger"
	"github.com/AthanatiusC/url-shortener/model"
	"github.com/AthanatiusC/url-shortener/repository"
)

type ShortenerService struct {
	config config.Config
	redis  *repository.ShortenerRepository
}

func NewShortenerService(ctx context.Context, config config.Config, shrtnRepository *repository.ShortenerRepository) *ShortenerService {
	rand.Seed(time.Now().UnixNano()) // initialize global pseudo random generator
	return &ShortenerService{
		config: config,
		redis:  shrtnRepository,
	}
}

func (s *ShortenerService) ShortenURL(ctx context.Context, URL string, exp string) (model.URLShortener, error) {
	if URL == "" {
		err := fmt.Errorf("URL cannot be empty")
		logger.ErrorContext(ctx, err)
		return model.URLShortener{HttpCode: http.StatusBadRequest}, err
	}

	url, err := url.ParseRequestURI(URL)
	if err != nil {
		err := fmt.Errorf("URL is not valid")
		logger.ErrorContext(ctx, err)
		return model.URLShortener{HttpCode: http.StatusBadRequest}, err
	}

	expiry, _ := time.ParseDuration(s.config.Application.DefaultExpiration)
	if exp != "" {
		expiry, err = time.ParseDuration(exp)
		if err != nil && exp != "" {
			err := fmt.Errorf("expiry format is not valid")
			logger.ErrorContext(ctx, err)
			return model.URLShortener{HttpCode: http.StatusBadRequest}, err
		}
	}

	generatedUrl := make([]rune, s.config.Application.GeneratedURLLength)
	for i := range generatedUrl {
		generatedUrl[i] = model.LetterRunes[rand.Intn(len(model.LetterRunes))]
	}

	data := model.URLShortener{
		URL:       url.String(),
		ShortURL:  string(generatedUrl),
		ExpiredAt: time.Now().Add(expiry).Format(time.RFC3339),
		Exp:       expiry,
	}

	err = s.redis.SetShortenURL(context.Background(), data)
	if err != nil {
		logger.ErrorContext(ctx, err)
		return model.URLShortener{}, err
	}
	return data, nil
}

func (s *ShortenerService) GetURL(ctx context.Context, shortURL string) (model.URLShortener, error) {
	if shortURL == "" {
		err := fmt.Errorf("shortURL cannot be empty")
		logger.ErrorContext(ctx, err)
		return model.URLShortener{HttpCode: http.StatusBadRequest}, err
	}

	if strings.Contains(shortURL, "http") {
		url, err := url.ParseRequestURI(shortURL)
		if err != nil {
			err := fmt.Errorf("URL is not valid")
			logger.ErrorContext(ctx, err)
			return model.URLShortener{HttpCode: http.StatusBadRequest}, err
		}
		shortURL = url.Path
	}

	response, err := s.redis.GetShortenURL(context.Background(), shortURL)
	if err != nil && err.Error() != "redis: nil" {
		logger.ErrorContext(ctx, err)
		return model.URLShortener{}, err
	}

	if response.ShortURL == "" {
		err := fmt.Errorf("url does not exist or has expired")
		logger.ErrorContext(ctx, err)
		return model.URLShortener{HttpCode: http.StatusNotFound}, err
	}

	return response, nil
}

func (s *ShortenerService) HandleShortenURL(ctx context.Context, shortURL string) (model.URLShortener, error) {
	shortenURL, err := s.redis.GetShortenURL(context.Background(), shortURL)
	if err != nil && err.Error() != "redis: nil" {
		logger.ErrorContext(ctx, err)
		return model.URLShortener{}, err
	}

	if shortenURL.ShortURL == "" {
		err := fmt.Errorf("url does not exist or has expired")
		logger.ErrorContext(ctx, err)
		return model.URLShortener{HttpCode: http.StatusNotFound}, err
	}

	shortenURL.ClickCount += 1
	err = s.redis.SetShortenURL(context.Background(), shortenURL)
	if err != nil {
		logger.ErrorContext(ctx, err)
		return model.URLShortener{}, err
	}

	return shortenURL, nil
}

func (s *ShortenerService) DeleteURL(ctx context.Context, shortURL string) (int, error) {
	shortenURL, err := s.redis.GetShortenURL(context.Background(), shortURL)
	if err != nil && err.Error() != "redis: nil" {
		logger.ErrorContext(ctx, err)
		return http.StatusNotFound, err
	}

	if shortenURL.ShortURL == "" {
		err := fmt.Errorf("url does not exist or has expired")
		logger.ErrorContext(ctx, err)
		return http.StatusNotFound, err
	}

	err = s.redis.DeleteShortenURL(ctx, shortenURL.ShortURL)
	if err != nil {
		logger.ErrorContext(ctx, err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
