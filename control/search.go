package control

import (
	"github.com/zmb3/spotify"
)

func searchTrack(q string) ([]spotify.FullTrack, error) {
	q = trackQuery(q)
	page, err := client.Search(q, spotify.SearchTypeTrack)
	if err != nil {
		return nil, err
	}

	if page.Tracks == nil {
		return nil, nil
	}
	return page.Tracks.Tracks, nil
}

func searchAlbum(q string) ([]spotify.SimpleAlbum, error) {
	q = albumQuery(q)
	page, err := client.Search(q, spotify.SearchTypeAlbum)
	if err != nil {
		return nil, err
	}
	if page.Albums == nil {
		return nil, nil
	}
	return page.Albums.Albums, nil
}

func searchArtist(q string) ([]spotify.FullArtist, error) {
	page, err := client.Search(q, spotify.SearchTypeArtist)
	if err != nil {
		return nil, err
	}
	if page.Artists == nil {
		return nil, nil
	}
	return page.Artists.Artists, nil
}

func searchPlaylist(q string) ([]Playlist, error) {
	page, err := client.Search(q, spotify.SearchTypePlaylist)
	if err != nil {
		return nil, err
	}
	if page.Playlists == nil {
		return nil, nil
	}
	pls := make([]Playlist, 0, len(page.Playlists.Playlists))
	for _, p := range page.Playlists.Playlists {
		pls = append(pls, plFromSimple(p))
	}

	return pls, nil
}
