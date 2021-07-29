package control

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/insomnimus/libman/handler/cmd"
	"github.com/zmb3/spotify"
)

func togglePlay() error {
	if err := updateDevice(); err != nil {
		return err
	}
	isPlaying = !isPlaying
	var err error
	if isPlaying {
		err = client.PlayOpt(&spotify.PlayOptions{
			DeviceID: &device.ID,
		})
	} else {
		err = client.PauseOpt(&spotify.PlayOptions{
			DeviceID: &device.ID,
		})
	}
	return err
}

func setVolume(n int) error {
	if err := updateDevice(); err != nil {
		return err
	}
	if n < 0 {
		n = 0
	} else if n > 100 {
		n = 100
	}
	err := client.VolumeOpt(n, &spotify.PlayOptions{
		DeviceID: &device.ID,
	})
	if err != nil {
		return err
	}
	device.Volume = n
	return nil
}

func adjustVolume(n int) error {
	if err := updateDevice(); err != nil {
		return err
	}
	n += device.Volume
	if n > 100 {
		n = 100
	} else if n < 0 {
		n = 0
	}
	err := client.VolumeOpt(n, &spotify.PlayOptions{
		DeviceID: &device.ID,
	})
	if err != nil {
		return err
	}
	device.Volume = n
	return nil
}

func getActiveDevice() (*spotify.PlayerDevice, error) {
	devs, err := client.PlayerDevices()
	if err != nil {
		return nil, err
	}
	for _, d := range devs {
		if d.Active && !d.Restricted {
			return &d, nil
		}
	}

	return nil, fmt.Errorf("no active device detected")
}

func updateDevice() error {
	if device == nil {
		dev, err := getActiveDevice()
		if err != nil {
			return err
		}
		device = dev
	}
	return nil
}

func playNext() error {
	if err := updateDevice(); err != nil {
		return err
	}
	err := client.NextOpt(&spotify.PlayOptions{
		DeviceID: &device.ID,
	})
	if err != nil {
		return err
	}
	isPlaying = true
	return nil
}

func playPrev() error {
	if err := updateDevice(); err != nil {
		return err
	}
	err := client.PreviousOpt(&spotify.PlayOptions{
		DeviceID: &device.ID,
	})
	if err != nil {
		return err
	}
	isPlaying = true
	return nil
}

func handleShuffle(arg string) error {
	if err := updateDevice(); err != nil {
		return err
	}
	var b bool
	switch strings.ToLower(arg) {
	case "on", "true", "yes":
		b = true
	case "off", "no", "false":
		b = false
	case "":
		b = !shuffleState
	default:
		handlers.ShowUsage(cmd.Shuffle)
		return nil
	}
	err := client.ShuffleOpt(b, &spotify.PlayOptions{
		DeviceID: &device.ID,
	})
	if err != nil {
		return err
	}
	if b {
		fmt.Println("shuffle on")
	} else {
		fmt.Println("shuffle off")
	}
	shuffleState = b
	return nil
}

func handleRepeat(arg string) error {
	if err := updateDevice(); err != nil {
		return err
	}
	var r string
	switch strings.ToLower(arg) {
	case "off", "false", "no":
		r = "off"
	case "on", "track", "song":
		r = "track"
	case "context":
		r = "context"
	default:
		handlers.ShowUsage(cmd.Repeat)
		return nil
	}

	err := client.RepeatOpt(r, &spotify.PlayOptions{
		DeviceID: &device.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("repeat = %s\n", r)
	repeatState = r
	return nil
}

func getPlaying() (*spotify.FullTrack, error) {
	cp, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		return nil, err
	}
	isPlaying = cp.Playing
	if cp.Item == nil {
		return nil, fmt.Errorf("not playing a track")
	}
	return cp.Item, nil
}

func handleVolume(arg string) error {
	if arg == "" {
		if device == nil {
			return fmt.Errorf("no active device detected")
		}
		fmt.Printf("The volume is %d%%.\n", device.Volume)
		return nil
	}

	n, err := strconv.Atoi(arg)
	if err != nil {
		handlers.ShowUsage(cmd.Volume)
		return nil
	}

	return setVolume(n)
}

func handleSetDevice(arg string) error {
	devs, err := client.PlayerDevices()
	if err != nil {
		return err
	}
	if len(devs) == 0 {
		fmt.Println("Couldn't detect any device.")
		return nil
	}

	var dev *spotify.PlayerDevice
	if arg != "" {
		for _, d := range devs {
			if strings.EqualFold(d.Name, arg) {
				dev = &d
				break
			}
		}

		if dev == nil {
			fmt.Printf("Did not find any device named %s.\n", arg)
			return nil
		}
	} else {
		for i, d := range devs {
			fmt.Printf("#%-2d | %s\n", i, d.Name)
		}
		n := readNumber(0, len(devs))
		if n < 0 {
			fmt.Println("cancelled")
			return nil
		}
		dev = &devs[n]
	}

	play := isPlaying

	err = client.TransferPlayback(dev.ID, play)
	if err != nil {
		return err
	}
	fmt.Printf("Playing on %s.\n", dev.Name)
	device = dev
	return nil
}

func handleShow(arg string) error {
	if err := updateDevice(); err != nil {
		return err
	}
	state, err := client.PlayerState()
	if err != nil {
		return err
	}
	device = &state.Device
	shuffleState = state.ShuffleState
	repeatState = state.RepeatState
	isPlaying = state.Playing
	t := state.Item
	if t == nil {
		fmt.Println("Not playing a track.")
	} else {
		fmt.Printf("Currently playing %s [%s] by %s.\n", t.Name, t.Album.Name, joinArtists(t.Artists))
		fmt.Printf("shuffle = %t\nrepeat = %s\n", shuffleState, repeatState)
	}
	return nil
}

func handleSharePlaying(arg string) error {
	if err := updateDevice(); err != nil {
		return err
	}
	state, err := client.PlayerState()
	if err != nil {
		return err
	}
	t := state.Item

	if t == nil {
		fmt.Println("Not playing a track.")
		return nil
	}
	err = clipboard.WriteAll(t.ExternalURLs["spotify"])
	if err != nil {
		return err
	}
	fmt.Printf("Copied the URL for %s [%s] by %s to the clipboard.\n", t.Name, t.Album.Name, joinArtists(t.Artists))
	return nil
}

func queueTrack(t *spotify.FullTrack) error {
	err := client.QueueSong(t.ID)
	if err == nil {
		fmt.Printf("Added %s by %s to the queue.\n", t.Name, joinArtists(t.Artists))
	}
	return err
}

func playUserLibrary() error {
	if err := updateLibraryCache(); err != nil {
		return err
	}
	if len(libraryCache) == 0 {
		fmt.Println("You have no saved tracks in your library.")
		return nil
	}
	err := client.PlayOpt(&spotify.PlayOptions{
		URIs: libraryCache.uris(),
	})
	if err == nil {
		fmt.Println("Playing tracks from your library.")
		isPlaying = true
	}
	return err
}
