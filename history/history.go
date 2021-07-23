package history

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var HistorySize uint32 = 100

type History struct {
	Artists  []string `json:"artists"`
	Albums   []string `json:"albums"`
	Tracks   []string `json:"tracks"`
	modified bool
}

func NewHistory() *History {
	return &History{
		Tracks:  make([]string, 0),
		Albums:  make([]string, 0),
		Artists: make([]string, 0),
	}
}

func (h *History) AppendArtist(s string) {
	h.modified = true
	h.Artists = append(h.Artists, s)
	if len(h.Artists) > int(HistorySize) {
		h.Artists = h.Artists[:HistorySize]
	}
}

func (h *History) AppendAlbum(s string) {
	h.modified = true
	h.Albums = append(h.Albums, s)
	if len(h.Albums) > int(HistorySize) {
		h.Albums = h.Albums[:HistorySize]
	}
}

func (h *History) AppendTrack(s string) {
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
