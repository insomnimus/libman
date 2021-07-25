package history

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

var HistorySize uint32 = 66

type History struct {
	Artists   []string `json:"artists"`
	Albums    []string `json:"albums"`
	Tracks    []string `json:"tracks"`
	Playlists []string `json:"playlists"`
	modified  bool
}

func NewHistory() *History {
	return &History{
		Tracks:    make([]string, 0),
		Albums:    make([]string, 0),
		Artists:   make([]string, 0),
		Playlists: make([]string, 0),
	}
}

func (h *History) AppendArtist(s string) {
	for _, a := range h.Artists {
		if strings.EqualFold(a, s) {
			return
		}
	}
	h.modified = true
	h.Artists = append(h.Artists, s)
	if len(h.Artists) > int(HistorySize) {
		h.Artists = h.Artists[:HistorySize]
	}
}

func (h *History) AppendAlbum(s string) {
	for _, a := range h.Albums {
		if strings.EqualFold(s, a) {
			return
		}
	}
	h.modified = true
	h.Albums = append(h.Albums, s)
	if len(h.Albums) > int(HistorySize) {
		h.Albums = h.Albums[:HistorySize]
	}
}

func (h *History) AppendTrack(s string) {
	for _, t := range h.Tracks {
		if strings.EqualFold(t, s) {
			return
		}
	}
	h.modified = true
	h.Tracks = append(h.Tracks, s)
	if len(h.Tracks) > int(HistorySize) {
		h.Tracks = h.Tracks[:HistorySize]
	}
}

func (h *History) Modified() bool {
	return h.modified
}

func Load(file string) (*History, error) {
	// init the file if it doesn't exist
	if _, e := os.Stat(file); errors.Is(e, os.ErrNotExist) {
		h := History{}
		err := h.Save(file)
		if err != nil {
			return nil, fmt.Errorf("error initializing the history file: %w", err)
		}
		return &h, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading the history file: %w", err)
	}
	var h History
	err = json.Unmarshal(data, &h)
	if err != nil {
		return nil, fmt.Errorf("error parsing the history file: %w", err)
	}
	return &h, nil
}

func (h *History) Save(file string) error {
	data, err := json.MarshalIndent(h, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o644)
}

func (h *History) AppendPlaylist(s string) {
	for _, p := range h.Playlists {
		if strings.EqualFold(p, s) {
			return
		}
	}
	h.modified = true
	h.Playlists = append(h.Playlists, s)
	if len(h.Playlists) > int(HistorySize) {
		h.Playlists = h.Playlists[:HistorySize]
	}
}
