package store

import (
	"encoding/json"
	"os"
	"time"

	"github.com/zmb3/spotify"
)

type Playlist struct {
	Name         string                  `json:"name"`
	ID           spotify.ID              `json:"id"`
	Owner        string                  `json:"owner"`
	DateExported time.Time               `json:"date_exported"`
	Tracks       []spotify.PlaylistTrack `json:"tracks"`
}

func PlaylistFromFull(p spotify.FullPlaylist) Playlist {
	return Playlist{
		Name:   p.Name,
		ID:     p.ID,
		Owner:  p.Owner.DisplayName,
		Tracks: p.Tracks.Tracks,
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
