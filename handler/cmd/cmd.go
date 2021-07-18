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
	Alias
	Share
)
