package set

import (
	"github.com/exceptioon/tiktok-fav-publisher/internal"
)

type setCache struct {
	cache map[string]struct{}
}

func NewSet() (internal.Cache, error) {
	s := setCache{
		cache: make(map[string]struct{}),
	}
	return &s, nil
}

func (s *setCache) Add(value string) error {
	s.cache[value] = struct{}{}
	return nil
}

func (s *setCache) IsExist(value string) bool {
	_, ok := s.cache[value]
	return ok
}
