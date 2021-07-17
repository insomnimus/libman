package alias

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Alias struct {
	Left      string
	Right     string
	timestamp time.Time
}

func (a Alias) String() string {
	return fmt.Sprintf("alias %s=%s", a.Left, a.Right)
}

type Set struct {
	m map[string]*Alias
}

func (s *Set) Set(left, right string) {
	if s.m == nil {
		s.m = make(map[string]*Alias)
	}
	s.m[strings.ToUpper(left)] = &Alias{
		Left:      left,
		Right:     right,
		timestamp: time.Now(),
	}
}

func (s *Set) Get(key string) (*Alias, bool) {
	if s.m == nil {
		return nil, false
	}
	val, ok := s.m[strings.ToUpper(key)]
	return val, ok
}

func (s *Set) Sorted() []*Alias {
	if s.m == nil {
		return nil
	}
	sl := make([]*Alias, 0, len(s.m))
	for _, a := range s.m {
		sl = append(sl, a)
	}

	sort.Slice(sl, func(i, j int) bool {
		return sl[i].timestamp.Before(sl[j].timestamp)
	})

	return sl
}

func (s *Set) Len() int {
	return len(s.m)
}

func (s *Set) Unset(key string) bool {
	if s.m == nil {
		return false
	}
	key = strings.ToUpper(key)
	_, ok := s.m[key]
	if ok {
		delete(s.m, key)
	}
	return ok
}

func (s *Set) Inner() map[string]*Alias {
	return s.m
}
