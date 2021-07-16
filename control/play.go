package control

import (
	"fmt"
	"github.com/zmb3/spotify"
)

func playTrack(t *spotify.FullTrack) error {
	err := client.PlayOpt(&spotify.PlayOptions{
		DeviceID: deviceID,
		URIs:     []spotify.URI{t.URI},
	})
	if err != nil {
		return err
	}

	fmt.Printf("Playing %s by %s.\n", t.Name, joinArtists(t.Artists))
	isPlaying = true
	return nil
}

func playAlbum(a *spotify.SimpleAlbum) error {
	err := client.PlayOpt(&spotify.PlayOptions{
		DeviceID:        deviceID,
		PlaybackContext: &a.URI,
		PlaybackOffset:  &spotify.PlaybackOffset{Position: 0},
	})
	if err != nil {
		return err
	}

	fmt.Printf("Playing %s by %s.\n", a.Name, joinArtists(a.Artists))
	isPlaying = true
	return nil
}

func playArtist(a *spotify.FullArtist) error {
	err := client.PlayOpt(&spotify.PlayOptions{
		DeviceID:        deviceID,
		PlaybackContext: &a.URI,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Playing tracks from %s.\n", a.Name)
	isPlaying = true
	return nil
}

func playPlaylist(p *Playlist) error {
	err := client.PlayOpt(&spotify.PlayOptions{
		DeviceID:        deviceID,
		PlaybackContext: &p.URI,
		PlaybackOffset:  &spotify.PlaybackOffset{Position: 0},
	})
	if err != nil {
		return err
	}

	fmt.Printf("Playing tracks from %s.\n", p.Name)
	isPlaying = true
	return nil
}
