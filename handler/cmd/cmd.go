package cmd

type Cmd uint8

//go:generate stringer -type=Cmd

const (
	_ Cmd = iota
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

	SavePlaying
	RemovePlaying
	CreatePlaylist
	EditPlaylist
	DeletePlaylist

	Help
	PlayUserPlaylist
	SetDevice
	Show
	Prompt
)
