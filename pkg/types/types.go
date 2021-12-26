package types

type Album struct {
	FlacAvailable bool
	Name          string
	MP3FileSize   string
	FlacFileSize  string
	FileCount     int
	Tracks        []Track
	Covers        []string
}

type TrackMetaIndexes struct {
	CDNumber     int
	TrackNumber  int
	SongName     int
	TrackLength  int
	MP3FileSize  int
	FlacFileSize int
}

type SearchResults map[string]string

type Track struct {
	CDNumber     string
	Number       string
	Name         string
	Duration     string
	MP3FileSize  string
	FlacFileSize string
	URL          string
}