package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AthanatiusC/url-shortener/model"
	"github.com/redis/go-redis/v9"
)

type ShortenerRepository struct {
	redis *redis.Client
}

func NewRedisRepository(ctx context.Context, redis *redis.Client) *ShortenerRepository {
	return &ShortenerRepository{
		redis: redis,
	}
}

const (
	URLKey = "url/%s"
)

func (r *ShortenerRepository) GetShortenURL(ctx context.Context, shortURL string) (model.URLShortener, error) {
	result, err := r.redis.Get(ctx, fmt.Sprintf(URLKey, shortURL)).Result()
	if err != nil {
		{
			return model.URLShortener{}, err
		}
	}
	var response model.URLShortener
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		return model.URLShortener{}, err
	}
	return response, nil
}

func (r *ShortenerRepository) SetShortenURL(ctx context.Context, request model.URLShortener) error {
	jsonString, err := json.Marshal(request)
	if err != nil {
		return err
	}
	return r.redis.Set(ctx, fmt.Sprintf(URLKey, request.ShortURL), string(jsonString), request.Exp).Err()
}

func (r *ShortenerRepository) DeleteShortenURL(ctx context.Context, shortURL string) error {
	return r.redis.Del(ctx, fmt.Sprintf(URLKey, shortURL)).Err()
}
