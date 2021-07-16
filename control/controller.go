package control

import (
	"github.com/zmb3/spotify"
	"libman/handler"
)

var (
	client       *spotify.Client
	user         *spotify.PrivateUser
	device       *spotify.PlayerDevice
	prompt       string
	handlers     handler.Set
	cache        = new(PlaylistCache)
	lastPl       *Playlist
	isPlaying    bool
	shuffleState bool
	repeatState  string
)
