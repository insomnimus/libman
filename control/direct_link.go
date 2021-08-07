package control

import (
	"fmt"
	"strings"

	"github.com/zmb3/spotify"
)

func handleLink(link string) error {
	link = strings.TrimPrefix(link, "https://open.spotify.com/")

	split := strings.SplitN(link, "/", 2)

	switch split[0] {
	case "track":
		return playTrackFromID(split[1])
	case "artist":
		return playArtistFromID(split[1])
	case "playlist":
		return playPlaylistFromID(split[1])
	case "album":
		return playAlbumFromID(split[1])
	default:
		return fmt.Errorf("unrecognized url")
	}
}

func playTrackFromID(id string) error {
	t, err := client.GetTrack(spotify.ID(id))
	if err != nil {
		return err
	}
	return playTrack(t)
}

func playArtistFromID(id string) error {
	a, err := client.GetArtist(spotify.ID(id))
	if err != nil {
		return err
	}
	return playArtist(a)
}

func playAlbumFromID(id string) error {
	a, err := client.GetAlbum(spotify.ID(id))
	if err != nil {
		return err
	}
	return playAlbum(&a.SimpleAlbum)
}

func playPlaylistFromID(id string) error {
	p, err := client.GetPlaylist(spotify.ID(id))
	if err != nil {
		return err
	}

	pl := Playlist{
		FullPlaylist: *p,
		isFull:       true,
		isFollowed:   false,
	}

	return playPlaylist(&pl)
}
