package model

import "time"

var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type ShortenerRequest struct {
	URL string `json:"url"`
	Exp string `json:"exp"`
}

type URLShortener struct {
	URL        string        `json:"url"`
	ShortURL   string        `json:"short_url"`
	ClickCount int           `json:"click_count"`
	ExpiredAt  string        `json:"expired_at"`
	Exp        time.Duration `json:"-"`
	HttpCode   int           `json:"-"`
}

type ShortenerResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}
