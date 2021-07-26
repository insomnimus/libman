package control

import (
	"fmt"
	"strings"

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
		// PlaybackOffset:  &spotify.PlaybackOffset{Position: 0}, // position being 0 makes json omit the field
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
	var p *Playlist
	if strings.Contains(arg, "::") {
		split := strings.SplitN(arg, "::", 2)
		left := strings.TrimSpace(split[0])
		var right string
		if len(split) > 1 {
			right = strings.TrimSpace(split[1])
		}
		p = choosePlaylist(left)
		if p == nil {
			return nil
		}
		if right != "" {
			t, err := p.findTrack(right)
			if err != nil {
				return err
			}
			// play with context as offset
			err = client.PlayOpt(&spotify.PlayOptions{
				PlaybackContext: &p.URI,
				PlaybackOffset: &spotify.PlaybackOffset{
					URI: t.URI,
				},
			})
			if err == nil {
				fmt.Printf("Playing %s from %s.\n", t.Name, p.Name)
				isPlaying = true
			}
			return err
		}
	}

	if p == nil {
		p = choosePlaylist(arg)
	}
	if p == nil {
		return nil
	}

	err := playPlaylist(p)
	if err != nil {
		return err
	}
	lastPl = p
	return nil
}
