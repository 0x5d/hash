package core

import (
	"context"
	"fmt"
	neturl "net/url"

	"go.uber.org/zap"
)

type URLService struct {
	urlRepo URLRepository
	cache   Cache
	log     *zap.Logger
}

type ShortenedURL struct {
	Original string
	ShortKey string
	Enabled  bool
}

type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "URL not found"
}

type ErrInvalidURL struct {
	URL string
}

func (e *ErrInvalidURL) Error() string {
	return fmt.Sprintf("invalid URL %q", e.URL)
}

func NewURLService(log *zap.Logger, urlRepo URLRepository, cache Cache) URLService {
	return URLService{urlRepo: urlRepo, cache: cache, log: log}
}

func (s *URLService) ShortenAndSave(ctx context.Context, url string, enabled bool) (*ShortenedURL, error) {
	_, err := neturl.Parse(url)
	if err != nil {
		return nil, &ErrInvalidURL{url}
	}
	id, err := s.urlRepo.NextID(ctx)
	if err != nil {
		return nil, err
	}
	str := base62Encode(id)
	short := ShortenedURL{Original: url, ShortKey: string(str), Enabled: enabled}
	err = s.urlRepo.Create(ctx, short)
	if err != nil {
		return nil, err
	}
	s.setAsync(string(str), url)
	return &short, nil
}

func (s *URLService) DecodeAndGet(ctx context.Context, key string) (*ShortenedURL, error) {
	url, err := s.cache.Get(ctx, key)
	if err != nil {
		s.log.Warn("error while fetching value from cache", zap.Error(err))
	}
	if url != "" {
		return &ShortenedURL{Original: url, ShortKey: key, Enabled: true}, nil
	}
	id, err := base62Decode([]byte(key))
	if err != nil {
		return nil, err
	}
	short, err := s.urlRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	s.setAsync(key, short.Original)
	return short, nil
}

func (s *URLService) Update(ctx context.Context, key string, newURL *string, enabled *bool) (*ShortenedURL, error) {
	id, err := base62Decode([]byte(key))
	if err != nil {
		return nil, err
	}
	short, err := s.urlRepo.Update(ctx, id, newURL, enabled)
	if err != nil {
		return nil, err
	}
	short.ShortKey = key
	if short.Enabled {
		s.setAsync(key, short.Original)
	} else {
		err = s.cache.Del(ctx, key)
		if err != nil {
			s.log.Warn("error while deleting key from cache", zap.Error(err))
		}
	}
	return short, nil
}

func (s *URLService) setAsync(key, url string) {
	go func() {
		err := s.cache.Set(context.Background(), key, url)
		if err != nil {
			s.log.Warn("error while setting key-value in cache", zap.Error(err))
		}
	}()
}
