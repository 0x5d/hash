package core

import (
	"context"
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

func NewURLService(log *zap.Logger, urlRepo URLRepository, cache Cache) URLService {
	return URLService{urlRepo: urlRepo, cache: cache, log: log}
}

func (s *URLService) ShortenAndSave(ctx context.Context, url string, enabled bool) (*ShortenedURL, error) {
	// TODO: create specific errors.
	_, err := neturl.Parse(url)
	if err != nil {
		return nil, err
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

func (s *URLService) setAsync(key, url string) {
	go func() {
		err := s.cache.Set(context.Background(), key, url)
		if err != nil {
			s.log.Warn("error while setting key-value in cache", zap.Error(err))
		}
	}()
}
