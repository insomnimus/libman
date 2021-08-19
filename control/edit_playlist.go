package control

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/insomnimus/libman/handler"
	"github.com/insomnimus/libman/handler/plcmd"
	"github.com/zmb3/spotify"
)

type PlaylistBuf struct {
	pl        spotify.FullPlaylist
	lastState []spotify.PlaylistTrack
	removes   []spotify.PlaylistTrack
	adds      []spotify.FullTrack

	page      int
	startFrom int
	upTo      int
	modified  bool
	handlers  handler.Set
}

func (p *Playlist) editTracks() error {
	if p.isFollowed {
		fmt.Println("You cannot edit this playlist.")
		return nil
	}
	if err := p.makeFull(); err != nil {
		return err
	}

	// create a temporary buffer containing the tracks
	buf := NewPlaylistBuf(p.FullPlaylist)
	modified, err := buf.interactive()
	if err != nil {
		return err
	}
	if modified {
		// update local cache entry
		pl, err := client.GetPlaylist(p.ID)
		if err != nil {
			return fmt.Errorf("error updating the local cache: %w", err)
		}
		*p = Playlist{*pl, true, false}
	}
	return nil
}

func (p *PlaylistBuf) add(track spotify.FullTrack) {
	// check if the track is already in the playlist
	for _, t := range p.pl.Tracks.Tracks {
		if t.Track.ID == track.ID {
			if readBool("%s is already in the playlist, do you still want to add it?", track.Name) {
				break
			}
			fmt.Println("Not added.")
			return
		}
	}

	// check the remove list
	for i, t := range p.removes {
		if t.Track.ID == track.ID {
			fmt.Printf("Removed %s from the remove queue.\n", track.Name)
			p.removes = append(p.removes[:i], p.removes[i+1:]...)
			p.pl.Tracks.Tracks = append(p.pl.Tracks.Tracks, t)
			return
		}
	}

	p.adds = append(p.adds, track)
	p.pl.Tracks.Tracks = append(p.pl.Tracks.Tracks, spotify.PlaylistTrack{Track: track})
	fmt.Printf("Added %s to the add queue.\n", track.Name)
	p.updateIndices()
}

func (p *PlaylistBuf) prevPage() {
	p.updateIndices()
	if p.page == 0 {
		fmt.Println("Already on the first page.")
		return
	}
	p.page -= 1

	p.upTo = p.startFrom
	p.startFrom -= PlaylistPageSize

	if p.upTo > len(p.pl.Tracks.Tracks) {
		p.upTo = len(p.pl.Tracks.Tracks)
	}

	p.displayPage()
}

func (p *PlaylistBuf) nextPage() {
	p.updateIndices()
	if p.startFrom+PlaylistPageSize > len(p.pl.Tracks.Tracks) {
		fmt.Println("Already on the last page.")
		return
	}

	p.page += 1
	p.startFrom += PlaylistPageSize
	p.upTo = p.startFrom + PlaylistPageSize
	if p.upTo > len(p.pl.Tracks.Tracks) {
		p.upTo = len(p.pl.Tracks.Tracks)
	}

	p.displayPage()
}

func (p *PlaylistBuf) displayPage() {
	p.updateIndices()
	for i, t := range p.pl.Tracks.Tracks[p.startFrom:p.upTo] {
		fmt.Printf("#%3d | %s by %s\n", i+p.startFrom, t.Track.Name, joinArtists(t.Track.Artists))
	}
}

func (p *PlaylistBuf) nPages() int {
	if len(p.pl.Tracks.Tracks)%PlaylistPageSize == 0 {
		return len(p.pl.Tracks.Tracks) / PlaylistPageSize
	}
	return len(p.pl.Tracks.Tracks)/PlaylistPageSize + 1
}

func (p *PlaylistBuf) interactive() (bool, error) {
	fmt.Printf("Editing %s.\nType `help` for a list of available commands.\n", p.pl.Name)
	fmt.Printf("Displaying page 1 of %d.\n", p.nPages())
	var input, cmd, arg string
	var cancelled bool
	p.displayPage()
	for {
		rl.SetCompleter(p.handlers.Complete)
		input, cancelled = readPrompt(true, "%s$ ", p.pl.Name)
		if cancelled {
			if p.hasChanges() {
				if readBool("You have unsaved changes, discard them?") {
					fmt.Println("Discarded all changes.")
					return p.modified, nil
				}
				continue
			}
			fmt.Println("returning")
			return p.modified, nil
		}
		if input == "" {
			if err := togglePlay(); err != nil {
				fmt.Printf("Error: %s.\n", err)
			}
			continue
		}
		cmd, arg = splitCmd(input)

		h := p.handlers.Match(cmd)
		if h == nil {
			fmt.Printf("%s is not a known command or alias.\nRun `help` for a list of available commands.\n", cmd)
			continue
		}
		switch h.Cmd {
		case plcmd.Return:
			if p.hasChanges() {
				if readBool("You have unsaved changes, discard them?") {
					fmt.Println("Discarded changes.")
					return p.modified, nil
				}
				continue
			}
			fmt.Println("returning")
			return p.modified, nil
		default:
			h.Run(arg)
		}
	}
}

func NewPlaylistBuf(pl spotify.FullPlaylist) *PlaylistBuf {
	upTo := len(pl.Tracks.Tracks)
	if upTo > PlaylistPageSize {
		upTo = PlaylistPageSize
	}
	lastState := make([]spotify.PlaylistTrack, len(pl.Tracks.Tracks))
	copy(lastState, pl.Tracks.Tracks)
	p := PlaylistBuf{
		pl:        pl,
		upTo:      upTo,
		lastState: lastState,
	}
	p.handlers = p.defaultHandlers()
	return &p
}

func (p *PlaylistBuf) defaultHandlers() handler.Set {
	hand := func(
		cmd uint8,
		name string,
		about string,
		usage string,
		desc string,
		aliases []string,
		run func(string),
	) handler.Handler {
		return handler.Handler{
			Cmd:      cmd,
			Name:     name,
			About:    about,
			Usage:    usage,
			Help:     desc,
			Aliases:  aliases,
			Complete: completeNothing,
			Run: func(arg string) error {
				run(arg)
				return nil
			},
		}
	}

	set := handler.Set{
		hand(
			plcmd.NextPage,
			"next",
			"Display the next page.",
			"next",
			"Display the next page",
			[]string{">"},
			func(string) { p.nextPage() },
		),
		hand(
			plcmd.PrevPage,
			"prev",
			"Display the previous page.",
			"prev",
			"Display the previous page.",
			[]string{"<"},
			func(string) { p.prevPage() },
		),
		hand(
			plcmd.Remove,
			"remove",
			"Remove a track from the playlist.",
			"remove <N...>",
			"Queue 1 or more tracks for removal.",
			[]string{"rm", "del"},
			p.handleRemove,
		),
		hand(
			plcmd.Play,
			"play",
			"Play a track from the playlist.",
			"play <N>",
			"Play a track from the playlist.",
			[]string{"pl"},
			p.handlePlay,
		),
		hand(
			plcmd.Add,
			"add",
			"Search for a track to add to the playlist.",
			"add <track>",
			"Search for a track to add to the playlist.",
			nil,
			p.handleAdd,
		),
		hand(
			plcmd.Discard,
			"discard",
			"Discard all your changes without returning.",
			"discard",
			"Discard all your changes without returning.",
			[]string{"abort"},
			func(string) {
				if !p.hasChanges() {
					fmt.Println("There are no unsaved changes.")
					return
				}
				if readBool("Are you sure you want to discard all the changes?") {
					p.discardChanges()
					fmt.Println("Discarded all the changes.")
				} else {
					fmt.Println("cancelled")
				}
			},
		),
		hand(
			plcmd.Display,
			"display",
			"Display the current page or your changes.",
			"display [changes|added|removed]",
			"Display the current page or your unsaved changes.",
			[]string{"show", "list", "ls"},
			p.handleDisplay,
		),
		hand(
			plcmd.Apply,
			"apply",
			"Apply your changes to the playlist.",
			"apply",
			"Apply all the changes to the playlist.",
			[]string{"commit"},
			func(string) {
				err := p.applyChanges()
				if err != nil {
					fmt.Printf("Error applying changes: %s.\n", err)
					return
				}
				p.modified = true
			},
		),
		hand(
			plcmd.Help,
			"help",
			"Show a lsit of available commands.",
			"help [command]",
			"Show help about a command or list available commands.",
			nil,
			p.handleHelp,
		),
		hand(
			plcmd.Return,
			"return",
			"Discard changes and return.",
			"return",
			"Discard all unapplied changes and return. You will be prompted if there are unsaved changes.",
			[]string{"done"},
			func(string) {},
		),
	}

	// apply completions for some commands
	set.Find(plcmd.Add).Complete = dynamicCompleteFunc(&Hist.Tracks, "add")
	set.Find(plcmd.Help).Complete = set.CompleteHelp
	set.Find(plcmd.Display).Complete = newWordCompleter("changes", "added", "removed")

	return set
}

func (p *PlaylistBuf) applyChanges() error {
	if !p.hasChanges() {
		return nil
	}
	if len(p.removes) > 0 {
		ids := make([]spotify.ID, len(p.removes))
		for i, t := range p.removes {
			ids[i] = t.Track.ID
		}
		_, err := client.RemoveTracksFromPlaylist(p.pl.ID, ids...)
		if err != nil {
			return err
		}
		fmt.Printf("Removed %d tracks from %s.\n", len(ids), p.pl.Name)
		p.modified = true
		p.removes = []spotify.PlaylistTrack{}
		p.lastState = p.pl.Tracks.Tracks
	}

	if len(p.adds) > 0 {
		ids := make([]spotify.ID, len(p.adds))
		for i, t := range p.adds {
			ids[i] = t.ID
		}
		_, err := client.AddTracksToPlaylist(p.pl.ID, ids...)
		if err != nil {
			return err
		}
		p.modified = true
		fmt.Printf("Added %d new tracks to %s.\n", len(ids), p.pl.Name)
		p.adds = []spotify.FullTrack{}
		p.lastState = p.pl.Tracks.Tracks
	}
	return nil
}

func (p *PlaylistBuf) discardChanges() {
	p.pl.Tracks.Tracks = p.lastState
	p.adds = []spotify.FullTrack{}
	p.removes = []spotify.PlaylistTrack{}
}

func (p *PlaylistBuf) handleRemove(arg string) {
	if arg == "" {
		p.handlers.ShowUsage(plcmd.Remove)
		return
	}

	split := strings.Fields(arg)
	nums := make([]int, 0, len(split))
	for _, s := range split {
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("'%s' is not a valid number.\n", s)
			continue
		}
		if n < p.startFrom || n >= p.upTo {
			fmt.Printf("Please enter a value between %d and %d.\n", p.startFrom, p.upTo)
			continue
		}
		nums = append(nums, n)
	}

	if len(nums) > 0 {
		p.remove(nums...)
	}
}

func (p *PlaylistBuf) updateIndices() {
	for p.page*PlaylistPageSize > len(p.pl.Tracks.Tracks) {
		p.page--
	}

	p.startFrom = p.page * PlaylistPageSize
	p.upTo = p.startFrom + PlaylistPageSize
	if p.upTo > len(p.pl.Tracks.Tracks) {
		p.upTo = len(p.pl.Tracks.Tracks)
	}
}

func (p *PlaylistBuf) handleAdd(arg string) {
	if arg == "" {
		p.handlers.ShowUsage(plcmd.Add)
		return
	}
	tracks, err := searchTrack(arg)
	if err != nil {
		fmt.Printf("Error: %s.\n", err)
		return
	}
	if len(tracks) == 0 {
		fmt.Printf("No result for %s.\n", arg)
		return
	}

	for i, t := range tracks {
		fmt.Printf("#%2d | %s by %s\n", i, t.Name, joinArtists(t.Artists))
	}

	n := readNumber(0, len(tracks))
	if n < 0 {
		fmt.Println("cancelled")
		return
	}

	p.add(tracks[n])
}

func (p *PlaylistBuf) handleDisplay(arg string) {
	switch strings.ToLower(arg) {
	case "":
		p.displayPage()
	case "add", "added":
		if len(p.adds) == 0 {
			fmt.Println("No tracks queued for addition.")
			return
		}
		fmt.Println("Tracks queued for addition:")
		for _, t := range p.adds {
			fmt.Printf("%s by %s\n", t.Name, joinArtists(t.Artists))
		}
	case "del", "removal", "rm":
		if len(p.removes) == 0 {
			fmt.Println("No tracks queued for removal.")
			return
		}
		fmt.Println("Tracks queued for removal:")
		for _, t := range p.removes {
			fmt.Printf("%s by %s\n", t.Track.Name, joinArtists(t.Track.Artists))
		}
	default:
		p.handlers.ShowUsage(plcmd.Display)
	}
}

func (p *PlaylistBuf) handlePlay(arg string) {
	if arg == "" {
		p.handlers.ShowUsage(plcmd.Play)
		return
	}

	n, err := strconv.Atoi(arg)
	if err != nil {
		p.handlers.ShowUsage(plcmd.Play)
		return
	}
	if n < p.startFrom || n >= p.upTo {
		fmt.Printf("Please enter a value between %d and %d.\n", p.startFrom, p.upTo)
		return
	}

	// play the playlist but with an index
	err = client.PlayOpt(&spotify.PlayOptions{
		PlaybackContext: &p.pl.URI,
		PlaybackOffset: &spotify.PlaybackOffset{
			Position: n,
			// URI:      p.pl.Tracks.Tracks[n].Track.URI,
		},
	})

	if err != nil {
		fmt.Printf("Error: %s.\n", err)
		return
	}

	fmt.Printf("Playing %s from %s.\n", p.pl.Tracks.Tracks[n].Track.Name, p.pl.Name)
}

func (p *PlaylistBuf) handleHelp(arg string) {
	if arg == "" {
		for _, h := range p.handlers {
			fmt.Println(h.String())
		}
		return
	}

	h := p.handlers.Match(arg)
	if h == nil {
		fmt.Printf("%s is not a known command or alias.\nRun `help` for a list f available commands.\n", arg)
		return
	}

	fmt.Println(h.GoString())
}

func (p *PlaylistBuf) hasChanges() bool {
	return len(p.adds) != 0 || len(p.removes) != 0
}

func (p *PlaylistBuf) remove(indices ...int) {
	sort.Ints(indices)
	remaining := make([]spotify.PlaylistTrack, 0, len(p.pl.Tracks.Tracks)-len(indices))
	fmt.Println("Queued for removal:")
	for i, t := range p.pl.Tracks.Tracks {
		if len(indices) == 0 {
			break
		}
		if i == indices[0] {
			indices = indices[1:]
			p.removes = append(p.removes, t)
			fmt.Printf("-  %s by %s\n", t.Track.Name, joinArtists(t.Track.Artists))
		} else {
			remaining = append(remaining, t)
		}
	}

	p.pl.Tracks.Tracks = remaining
}
