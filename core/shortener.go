package core

import "fmt"

type Shortener struct {
	c  *Config
	ID uint64
}

type ShortenedURL struct {
	URL string
}

func NewShortener(c Config) Shortener {
	return Shortener{c: &c, ID: 0}
}

func (s *Shortener) Shorten(url string) (*ShortenedURL, error) {
	str := base62Encode(s.ID)
	s.ID++
	return &ShortenedURL{URL: fmt.Sprintf("%s/%s", s.c.AdvertisedAddr, str)}, nil
}
