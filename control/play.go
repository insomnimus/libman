package control

import (
	"fmt"

	"github.com/zmb3/spotify"
)

func playTrack(t *spotify.FullTrack) error {
	if err := updateDevice(); err != nil {
		return err
	}
	err := client.PlayOpt(&spotify.PlayOptions{
		DeviceID: &device.ID,
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
	if err := updateDevice(); err != nil {
		return err
	}
	err := client.PlayOpt(&spotify.PlayOptions{
		DeviceID:        &device.ID,
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
	if err := updateDevice(); err != nil {
		return err
	}
	err := client.PlayOpt(&spotify.PlayOptions{
		DeviceID:        &device.ID,
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
	if err := updateDevice(); err != nil {
		return err
	}
	err := client.PlayOpt(&spotify.PlayOptions{
		DeviceID:        &device.ID,
		PlaybackContext: &p.URI,
		// PlaybackOffset:  &spotify.PlaybackOffset{Position: 0},
	})
	if err != nil {
		return err
	}

	fmt.Printf("Playing tracks from %s.\n", p.Name)
	isPlaying = true
	return nil
}

func handlePlayUserPlaylist(arg string) error {
	pl := choosePlaylist(arg)
	if pl == nil {
		return nil
	}

	err := playPlaylist(pl)
	if err != nil {
		return err
	}
	lastPl = pl
	return nil
}
