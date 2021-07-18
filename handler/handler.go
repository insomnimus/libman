package handler

import (
	"fmt"
	"strings"
)

type Handler struct {
	Name     string
	Aliases  []string
	Cmd      uint8
	Help     string
	About    string
	Usage    string
	Run      func(string) error
	Complete func(string) []string
}

type Set []Handler

func (s *Set) Find(c uint8) *Handler {
	for i, h := range *s {
		if h.Cmd == c {
			return &(*s)[i]
		}
	}
	return nil
}

func (h Handler) String() string {
	if len(h.Aliases) == 0 {
		return fmt.Sprintf("#%s\n  %s", h.Name, h.About)
	}
	return fmt.Sprintf("#%s [aliases: %s]\n  %s", h.Name, strings.Join(h.Aliases, ", "), h.About)
}

func (h Handler) GoString() string {
	if len(h.Aliases) == 0 {
		return fmt.Sprintf("#%s\nusage:\n  %s\n  %s", h.Name, h.Usage, h.Help)
	}

	return fmt.Sprintf(`#%s
aliases:
  %s
usage:
  %s
  %s`, h.Name, strings.Join(h.Aliases, ", "), h.Usage, h.Help)
}

func (s Set) ShowUsage(c uint8) {
	if h := s.Find(c); h != nil {
		fmt.Printf("usage:\n  %s\n", h.Usage)
	}
}

func (s *Set) Match(cmd string) *Handler {
	for i, h := range *s {
		if h.Matches(cmd) {
			return &(*s)[i]
		}
	}
	return nil
}

func (h *Handler) Matches(s string) bool {
	if strings.EqualFold(h.Name, s) {
		return true
	}
	for _, a := range h.Aliases {
		if strings.EqualFold(a, s) {
			return true
		}
	}
	return false
}

func (s *Set) Len() int {
	return len(*s)
}

func (h *Handler) HasPrefix(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}
	s = strings.ToLower(s)
	if strings.HasPrefix(h.Name, s) {
		return true
	}

	for _, a := range h.Aliases {
		if strings.HasPrefix(s, a) {
			return true
		}
	}

	return false
}
