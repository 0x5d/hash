package core

import (
	"context"
	neturl "net/url"
)

type URLService struct {
	urlRepo URLRepository
}

type ShortenedURL struct {
	Original string
	ShortKey string
	Enabled  bool
}

func NewURLService(urlRepo URLRepository) URLService {
	return URLService{urlRepo: urlRepo}
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
	return &short, nil
}
