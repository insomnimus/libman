package control

import (
	"github.com/zmb3/spotify"
	"libman/handler"
)

var (
	client    *spotify.Client
	user      *spotify.PrivateUser
	prompt    string
	handlers  handler.Set
	cache     = new(PlaylistCache)
	lastPl    *Playlist
	isPlaying bool
	deviceID  *spotify.ID
)
