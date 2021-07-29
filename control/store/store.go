package store

import (
	"encoding/json"
	"os"
	"time"

	"github.com/zmb3/spotify"
)

type Playlist struct {
	Name         string      `json:"name"`
	URI          spotify.URI `json:"uri"`
	ID           spotify.ID  `json:"id"`
	Owner        string      `json:"owner"`
	DateExported time.Time   `json:"date_exported"`
	Tracks       []Track     `json:"tracks"`
}

type Track struct {
	Name    string      `json:"name"`
	ID      spotify.ID  `json:"id"`
	URI     spotify.URI `json:"uri"`
	Artists []Artist    `json:"artists"`
}

type Artist struct {
	Name string `json:"name"`

	ID  spotify.ID  `json:"id"`
	URI spotify.URI `json:"uri"`
}

func TrackFromFull(t spotify.FullTrack) Track {
	artists := make([]Artist, len(t.Artists))
	for i, a := range t.Artists {
		artists[i] = Artist{
			Name: a.Name,
			URI:  a.URI,
			ID:   a.ID,
		}
	}
	return Track{
		Name:    t.Name,
		Artists: artists,
		ID:      t.ID,
		URI:     t.URI,
	}
}

func (t Track) Full() spotify.FullTrack {
	artists := make([]spotify.SimpleArtist, len(t.Artists))
	for i, a := range t.Artists {
		artists[i] = a.Simple()
	}
	return spotify.FullTrack{SimpleTrack: spotify.SimpleTrack{
		Name:    t.Name,
		ID:      t.ID,
		URI:     t.URI,
		Artists: artists,
	}}
}

func (a Artist) Simple() spotify.SimpleArtist {
	return spotify.SimpleArtist{
		Name: a.Name,
		ID:   a.ID,
		URI:  a.URI,
	}
}

func PlaylistFromFull(p spotify.FullPlaylist) Playlist {
	tracks := make([]Track, len(p.Tracks.Tracks))
	for i, t := range p.Tracks.Tracks {
		tracks[i] = TrackFromFull(t.Track)
	}

	return Playlist{
		Name:   p.Name,
		URI:    p.URI,
		ID:     p.ID,
		Owner:  p.Owner.DisplayName,
		Tracks: tracks,
	}
}

func (p *Playlist) ExportTo(filename string) error {
	p.DateExported = time.Now()
	data, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0o764)
}

func (p Playlist) PlaylistTracks() []spotify.PlaylistTrack {
	tracks := make([]spotify.PlaylistTrack, len(p.Tracks))
	for i, t := range p.Tracks {
		tracks[i] = spotify.PlaylistTrack{
			Track: t.Full(),
		}
	}
	return tracks
}

func ImportFrom(filename string) (*Playlist, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var p Playlist
	err = json.Unmarshal(data, &p)
	return &p, err
}

func ExportTo(p spotify.FullPlaylist, filename string) error {
	pl := PlaylistFromFull(p)
	return pl.ExportTo(filename)
}
