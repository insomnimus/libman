package control

import (
	"fmt"
	"github.com/zmb3/spotify"
)

func (p *Playlist) makeFull() error {
	if p.isFull {
		return nil
	}

	pl, err := client.GetPlaylist(p.ID)
	if err != nil {
		return err
	}

	*p = Playlist{*pl, true}
	return nil
}

func (p *Playlist) addTrack(track spotify.FullTrack) error {
	if !p.isFull {
		err := p.makeFull()
		if err != nil {
			return err
		}
	}

	// do not add if the track is a duplicate
	for _, t := range p.Tracks.Tracks {
		if t.Track.ID == track.ID {
			fmt.Printf("%s is already in %s, no action taken.\n", track.Name, p.Name)
			return nil
		}
	}

	_, err := client.AddTracksToPlaylist(p.ID, track.ID)
	if err != nil {
		return err
	}

	p.Tracks.Tracks = append(p.Tracks.Tracks, spotify.PlaylistTrack{
		Track: track,
	})
	fmt.Printf("Added %s to %s.\n", track.Name, p.Name)
	return nil
}

func (p *Playlist) removeTrack(track spotify.FullTrack) error {
	if !p.isFull {
		err := p.makeFull()
		if err != nil {
			return err
		}
	}

	// do not send a request if the playlist doesn't have the track
	index := -1
	for i, t := range p.Tracks.Tracks {
		if t.Track.ID == track.ID {
			index = i
			break
		}
	}

	if index < 0 {
		fmt.Printf("%s is not in %s.\n", track.Name, p.Name)
		return nil
	}

	_, err := client.RemoveTracksFromPlaylist(p.ID, track.ID)
	if err != nil {
		return err
	}

	fmt.Printf("Removed %s from %s.\n", track.Name, p.Name)
	p.Tracks.Tracks = append(p.Tracks.Tracks[:index], p.Tracks.Tracks[index+1:]...)
	return nil
}

func (p *Playlist) editDetails() error {
	name := readString("playlist name (%s): ", p.Name)
	if name == "" {
		name = p.Name
	}
	fmt.Println(name)

	desc := readString("playlist description: ")

	if name != p.Name {
		err := client.ChangePlaylistName(p.ID, name)
		if err != nil {
			return err
		}
		fmt.Printf("Renamed %s -> %s.\n", p.Name, name)
	}

	if desc != "" {
		err := client.ChangePlaylistDescription(p.ID, desc)
		if err != nil {
			return err
		}
		fmt.Printf("Changed the description of %s.\n", p.Name)
	}

	return nil
}

func (p *Playlist) play() error {
	err := playPlaylist(p)
	if err != nil {
		return err
	}
	lastPl = p
	return nil
}
