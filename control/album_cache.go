package control

import "github.com/zmb3/spotify"

type AlbumCache []spotify.SimpleAlbum

func (c AlbumCache) contains(id spotify.ID) bool {
	for _, a := range c {
		if a.ID == id {
			return true
		}
	}
	return false
}

func (c *AlbumCache) push(a spotify.SimpleAlbum) {
	*c = append(*c, a)
}

func (c AlbumCache) get(n int) *spotify.SimpleAlbum {
	return &c[n]
}

func (c *AlbumCache) remove(n int) {
	if n+1 == len(*c) {
		*c = (*c)[:n]
		return
	}
	*c = append(
		(*c)[:n],
		(*c)[n+1:]...,
	)
}
