package control

import (
	"fmt"
	"sort"
	"strings"

	"github.com/insomnimus/libman/handler"
	"github.com/insomnimus/libman/handler/cmd"
)

func DefaultHandlers() handler.Set {
	hand := func(c uint8, name, about, usage, help string, aliases []string, fn func(string) error) handler.Handler {
		return handler.Handler{
			Cmd:      c,
			Name:     name,
			About:    about,
			Usage:    usage,
			Help:     help,
			Aliases:  aliases,
			Run:      fn,
			Complete: completeNothing,
		}
	}

	set := handler.Set{
		// library commands
		hand(
			cmd.EditPlaylistDetails,
			"edit-details",
			"Edit a playlists name or description.",
			"edit-details [playlist]",
			"Edit a playlists name or description.",
			[]string{"rename"},
			handleEditPlaylistDetails,
		),
		hand(
			cmd.LikePlaying,
			"like-playing",
			"Save the currently playing track to your 'my music' folder.",
			"save-playing",
			"Save the currently playing track to your 'my music' folder.",
			[]string{"like", "fav", "fave"},
			handleLikePlaying,
		),
		hand(
			cmd.DislikePlaying,
			"dislike-playing",
			"Remove the currently playing track from your 'my music' folder.",
			"dislike-playing",
			"Remove the currently playing track from your 'my music' folder.",
			[]string{"dislike", "unfave", "unfav"},
			handleDislikePlaying,
		),
		hand(
			cmd.SavePlaying,
			"save-playing",
			"Save the currently playing song to a playlist.",
			"save-playing [playlist name]",
			"Save the currently playing track to a playlist.",
			[]string{"save", "add"},
			handleSavePlaying,
		),
		hand(
			cmd.RemovePlaying,
			"remove-playing",
			"Remove the currently playing track from a playlist.",
			"remove-playing [playlist name]",
			"Remove the currently playing track from one of your playlists.\nIf the playlist name is not given, the last played playlist will be assumed.",
			[]string{"rm"},
			handleRemovePlaying,
		),
		hand(
			cmd.CreatePlaylist,
			"create-playlist",
			"Create a new playlist.",
			"create-playlist [name]",
			"Create a new playlist.\nYou will be prompted for the details.",
			[]string{"create", "new-pl"},
			handleCreatePlaylist,
		),
		hand(
			cmd.EditPlaylist,
			"edit-playlist",
			"Edit a playlist.",
			"edit-playlist [name]",
			"Edit one of your playlists.",
			[]string{"edit"},
			handleEditPlaylist,
		),
		hand(
			cmd.DeletePlaylist,
			"delete-playlist",
			"Delete a playlist.",
			"delete-playlist [name]",
			"Delete one of your playlists or unfollow a playlist you're following.",
			[]string{"del"},
			handleDeletePlaylist,
		),

		// search commands
		hand(
			cmd.SearchTrack,
			"search-track",
			"Search for a track.",
			"search-track <track>",
			"Search for a track.\nYou can use `track::artist` or `track by artist` to get songs from an artist.",
			[]string{"stra"},
			handleSTrack,
		),
		hand(
			cmd.SearchAlbum,
			"search-album",
			"Search for an album.",
			"search-album <album>",
			"Search for an album.\nYou can use `album::artist` or `album by artist` to get the albums by an artist.",
			[]string{"salb"},
			handleSAlbum,
		),
		hand(
			cmd.SearchArtist,
			"search-artist",
			"Search for an artist.",
			"search-artist <artist>",
			"Search for an artist.",
			[]string{"sart"},
			handleSArtist,
		),
		hand(
			cmd.SearchPlaylist,
			"search-playlist",
			"Search for a playlist.",
			"search-playlist <playlist>",
			"Search for a public playlist.",
			[]string{"spla"},
			handleSPlaylist,
		),

		// play-first commands
		hand(
			cmd.PlayFirstTrack,
			"play-track",
			"Search for a track and play the first result.",
			"play-track <track>",
			"Search for a track and play the first result.\nYou can use `track::artist` or `track by artist` to limit the search to a specific artist.",
			[]string{"ptra"},
			handlePFTrack,
		),
		hand(
			cmd.PlayFirstAlbum,
			"play-album",
			"Search for an album and play the first result.",
			"play-album <album>",
			"Search for an album and play the first result.\nYou can use `album::artist` or `album by artist` to limit the search to a specific artist.",
			[]string{"palb"},
			handlePFAlbum,
		),
		hand(
			cmd.PlayFirstArtist,
			"play-artist",
			"Search for an artist and play the first result.",
			"play-artist <artist>",
			"Search for an artist and play the first result.",
			[]string{"part"},
			handlePFArtist,
		),
		hand(
			cmd.PlayFirstPlaylist,
			"play-playlist",
			"Search for a playlist and play the first result.",
			"play-playlist <playlist>",
			"Search for a public playlist and play the first result.",
			[]string{"ppla"},
			handlePFPlaylist,
		),

		// player commands
		hand(
			cmd.Prev,
			"prev",
			"Play the previous track.",
			"prev",
			"Play the previous track.",
			[]string{"<"},
			func(string) error {
				return playPrev()
			},
		),
		hand(
			cmd.Next,
			"next",
			"Play the next track.",
			"next",
			"Play the next track.",
			[]string{">"},
			func(string) error {
				return playNext()
			},
		),
		hand(
			cmd.Volume,
			"volume",
			"Display or change the volume.",
			"volume [percentage]",
			"Display or set the volume.\nYou can also use the `+N` or `-N` commands to change the volume by N%.",
			[]string{"vol"},
			handleVolume,
		),
		hand(
			cmd.Shuffle,
			"shuffle",
			"Change the shuffle state.",
			"shuffle [on|off]",
			"Change the shuffle state.\nIf none no arguments are supplied, this will switch the current shuffle state.",
			[]string{"sh"},
			handleShuffle,
		),
		hand(
			cmd.Repeat,
			"repeat",
			"Change the repeat state.",
			"repeat <off|track|context>",
			"Change the repeat state.\nno: do not repeat\ntrack: repeat the current track\ncontext: repeat the current album/artist/playlist",
			[]string{"rep"},
			handleRepeat,
		),

		// misc commands
		hand(
			cmd.RelatedArtists,
			"related-artists",
			"Get a list of related artists for an artist.",
			"related-artists <artist>",
			"Get related artists for a given artist.",
			[]string{"related", "rel"},
			handleRelatedArtists,
		),
		hand(
			cmd.Recommend,
			"recommend",
			"Recommend tracks based on a playlist.",
			"recommend [playlist]::[engine]",
			`Generate recommendations based on a playlist.
Recommendation style can be altered by specifying one of the following as the engine:
-	normal: This is the default engine
-	extreme: This will focus on the more extreme attributes of a playlist.
-	min: This will choose a set of attributes and will request recommendations that are higher than its values.
-	max: the opposite of min

Example usage (assuming you have a playlist named "workout"):
recommend workout::extreme`,
			[]string{"rec"},
			handleRecommend,
		),
		hand(
			cmd.SetDevice,
			"device",
			"Change the playback device.",
			"device [name]",
			"Change the playback device.\nIf no name is given, you'll be prompted to choose a device from a list of available devices.",
			[]string{"dev"},
			handleSetDevice,
		),
		hand(
			cmd.Prompt,
			"prompt",
			"Change the libman prompt.",
			"prompt <new prompt>",
			"Change the libman prompt.\nA space character will automatically be added to the end.",
			nil,
			handlePrompt,
		),
		hand(
			cmd.PlayUserPlaylist,
			"play",
			"Play a playlist from your library.",
			"play [name]",
			"Play a playlist from your library.\nIf the playlist name is not given, you will be prompted to choose one from a list of your playlists.",
			[]string{"pl"},
			handlePlayUserPlaylist,
		),
		hand(
			cmd.PlayLibrary,
			"play-library",
			"Play tracks from your 'my music' folder.",
			"play-library",
			"Play tracks from your 'my library' folder.",
			[]string{"lib"},
			func(string) error {
				return playUserLibrary()
			},
		),
		hand(
			cmd.Help,
			"help",
			"Display help about a command.",
			"help [command or alias]",
			"Display help about a command or list available commands.",
			nil,
			handleHelp,
		),
		hand(
			cmd.Show,
			"show",
			"Show the currently playing track.",
			"show",
			"Show the currently playing track.",
			[]string{"sw"},
			handleShow,
		),
		hand(
			cmd.Alias,
			"alias",
			"Define aliases to commands.",
			"alias <alias>=<command>",
			"Define aliases to commands.",
			nil,
			handleAlias,
		),
		hand(
			cmd.Share,
			"share-playing",
			"Copy the link to the currently playing track.",
			"share-playing",
			"Copy the link to the currently playing track to your clipboard.",
			[]string{"share", "cp"},
			handleSharePlaying,
		),
	}

	_applySuggestPlaylist(set)
	_applySuggestShuffleAndRepeat(set)
	_applySuggestHelp(set)
	_applySuggestHistory(set)
	_applySuggestRecommend(set)

	return set
}

func _applySuggestPlaylist(set handler.Set) {
	set.Find(cmd.EditPlaylistDetails).Complete = suggestPlaylist
	set.Find(cmd.PlayUserPlaylist).Complete = suggestPlaylist
	set.Find(cmd.SavePlaying).Complete = suggestPlaylist
	set.Find(cmd.RemovePlaying).Complete = suggestPlaylist
	set.Find(cmd.EditPlaylist).Complete = suggestPlaylist
	set.Find(cmd.DeletePlaylist).Complete = suggestPlaylist
}

func _applySuggestShuffleAndRepeat(set handler.Set) {
	set.Find(cmd.Shuffle).Complete = newWordCompleter([]string{"on", "off"}, "sh", "shuffle")
	set.Find(cmd.Repeat).Complete = newWordCompleter([]string{"off", "track", "context"}, "rep", "repeat")
}

func _applySuggestHelp(set handler.Set) {
	topics := make([]string, 0, len(set))
	for _, h := range set {
		topics = append(topics, h.Name)
		topics = append(topics, h.Aliases...)
	}
	sort.Strings(topics)
	set.Find(cmd.Help).Complete = newWordCompleter(topics, "help")
}

func _applySuggestHistory(set handler.Set) {
	set.Find(cmd.PlayFirstArtist).Complete = dynamicCompleteFunc(&Hist.Artists, "play-artist", "part")
	set.Find(cmd.SearchArtist).Complete = dynamicCompleteFunc(&Hist.Artists, "search-artist", "sart")
	set.Find(cmd.RelatedArtists).Complete = dynamicCompleteFunc(&Hist.Artists, "related-artists", "related", "rel")

	set.Find(cmd.PlayFirstAlbum).Complete = dynamicCompleteFunc(&Hist.Albums, "palb", "play-album")
	set.Find(cmd.SearchAlbum).Complete = dynamicCompleteFunc(&Hist.Albums, "salb", "search-album")

	set.Find(cmd.PlayFirstTrack).Complete = dynamicCompleteFunc(&Hist.Tracks, "ptra", "play-track")
	set.Find(cmd.SearchTrack).Complete = dynamicCompleteFunc(&Hist.Tracks, "stra", "search-track")

	set.Find(cmd.PlayFirstPlaylist).Complete = dynamicCompleteFunc(&Hist.Playlists, "play-playlist", "ppla")
	set.Find(cmd.SearchPlaylist).Complete = dynamicCompleteFunc(&Hist.Playlists, "spla", "search-playlist")
}

func _applySuggestRecommend(set handler.Set) {
	set.Find(cmd.Recommend).Complete = func(buf string) []string {
		if err := updateCache(); err != nil {
			return nil
		}
		buf = strings.TrimPrefix(buf, " ")
		command, arg := splitCmd(buf)
		if !strings.EqualFold(command, "recommend") && !strings.EqualFold(command, "rec") {
			return nil
		}
		// complete playlist name if there's no ::
		c := make([]string, 0, len(cache))
		if !strings.Contains(arg, "::") {
			for _, p := range cache {
				if hasPrefixFold(p.Name, arg) {
					c = append(c, fmt.Sprintf("%s %s", command, p.Name))
				}
			}
		} else {
			// complete the engine
			split := strings.SplitN(arg, "::", 2)
			pl := split[0]
			arg = ""
			if len(split) > 1 {
				arg = split[1]
			}
			if hasPrefixFold("normal", arg) {
				c = append(c, fmt.Sprintf("%s %s::normal", command, pl))
			}
			if hasPrefixFold("extreme", arg) {
				c = append(c, fmt.Sprintf("%s %s::extreme", command, pl))
			}
			if hasPrefixFold("max", arg) {
				c = append(c, fmt.Sprintf("%s %s::max", command, pl))
			}
			if hasPrefixFold("min", arg) {
				c = append(c, fmt.Sprintf("%s %s::min", command, pl))
			}
		}

		if len(c) == 0 {
			return nil
		}
		return c
	}
}
