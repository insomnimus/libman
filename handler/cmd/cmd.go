package cmd

const (
	_ uint8 = iota
	Search
	SearchTrack
	SearchArtist
	SearchAlbum
	SearchPlaylist

	PlayFirstTrack
	PlayFirstAlbum
	PlayFirstArtist
	PlayFirstPlaylist

	Volume
	Shuffle
	Repeat
	Next
	Prev

	DislikePlaying
	LikePlaying
	SavePlaying
	RemovePlaying
	CreatePlaylist
	EditPlaylistDetails
	EditPlaylist
	DeletePlaylist

	Help
	PlayUserPlaylist
	PlayLibrary
	PlayTopTracks
	SetDevice
	Show
	Prompt
	Alias
	Share
	RelatedArtists
	Recommend
	ImportPlaylist
	ExportPlaylist
)
