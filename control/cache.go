package control

import (
	"strings"

	"github.com/zmb3/spotify"
)

type Playlist struct {
	spotify.FullPlaylist
	isFull     bool
	isFollowed bool
}

func plFromSimple(p spotify.SimplePlaylist) Playlist {
	return Playlist{
		spotify.FullPlaylist{SimplePlaylist: p},
		false,
		false,
	}
}

type PlaylistCache []Playlist

func (c *PlaylistCache) insertFull(index int, p spotify.FullPlaylist) {
	left := append((*c)[:index], Playlist{p, true, false})
	right := (*c)[index:]
	*c = append(left, right...)
}

func (c *PlaylistCache) push(p Playlist) {
	*c = append(PlaylistCache{p}, (*c)...)
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

func (c *PlaylistCache) findByName(s string) *Playlist {
	for i := range *c {
		if strings.EqualFold((*c)[i].Name, s) {
			return &(*c)[i]
		}
	}
	return nil
}

func (c *PlaylistCache) get(n int) *Playlist {
	return &(cache)[n]
}

func (c *PlaylistCache) pushFollowed(p spotify.SimplePlaylist) {
	pl := Playlist{
		FullPlaylist: spotify.FullPlaylist{SimplePlaylist: p},
		isFollowed:   true,
	}
	c.push(pl)
}
