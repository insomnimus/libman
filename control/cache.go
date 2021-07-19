package control

import (
	"strings"

	"github.com/zmb3/spotify"
)

type Playlist struct {
	spotify.FullPlaylist
	isFull bool
}

func plFromSimple(p spotify.SimplePlaylist) Playlist {
	return Playlist{
		spotify.FullPlaylist{SimplePlaylist: p},
		false,
	}
}

type PlaylistCache []Playlist

func (c *PlaylistCache) insertFull(index int, p spotify.FullPlaylist) {
	left := append((*c)[:index], Playlist{p, true})
	right := (*c)[index:]
	*c = append(left, right...)
}

func (c *PlaylistCache) insertSimple(index int, p spotify.SimplePlaylist) {
	left := append((*c)[:index], Playlist{
		spotify.FullPlaylist{SimplePlaylist: p}, false,
	})
	right := (*c)[index:]
	*c = append(left, right...)
}

func (c *PlaylistCache) remove(id spotify.ID) {
	index := -1
	for i, p := range *c {
		if p.ID == id {
			index = i
			break
		}
	}

	if index < 0 {
		return
	}

	if index+1 == len(*c) {
		*c = (*c)[:index]
		return
	}

	left := (*c)[:index]
	right := (*c)[index+1:]
	*c = append(left, right...)
}

func (c *PlaylistCache) pushSimple(p spotify.SimplePlaylist) {
	*c = append(*c, plFromSimple(p))
}

func (c *PlaylistCache) pushFull(p spotify.FullPlaylist) {
	*c = append(*c, Playlist{p, true})
}

func (c *PlaylistCache) find(id spotify.ID) *Playlist {
	for i := range *c {
		if (*c)[i].ID == id {
			return &(*c)[i]
		}
	}
	return nil
}

func (c *PlaylistCache) findByName(s string) *Playlist {
	for i := range *c {
		if strings.EqualFold((*c)[i].Name, s) {
			return &(*c)[i]
		}
	}
	return nil
}

func (c *PlaylistCache) get(n int) *Playlist {
	return &(*cache)[n]
}
