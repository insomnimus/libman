package control

import (
	"github.com/zmb3/spotify"
)

type LibraryCache []spotify.FullTrack

func (c *LibraryCache) contains(id spotify.ID) bool {
	for _, t := range *c {
		if t.ID == id {
			return true
		}
	}
	return false
}

func (c *LibraryCache) push(t spotify.FullTrack) {
	*c = append(*c, t)
}

func (c *LibraryCache) remove(n int) {
	if n+1 == len(*c) {
		*c = (*c)[:n]
		return
	}
	*c = append(
		(*c)[:n],
		(*c)[n+1:]...,
	)
}

func (c *LibraryCache) uris() (uris []spotify.URI) {
	for _, t := range *c {
		uris = append(uris, t.URI)
	}
	return
}

func (c *LibraryCache) removeByID(id spotify.ID) {
	index := -1
	for i, t := range *c {
		if t.ID == id {
			index = i
			break
		}
	}
	if index > 0 {
		c.remove(index)
	}
}
