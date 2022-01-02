package types

type Album struct {
	Slug          string
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

type SearchIndex struct {
	IndexVersion string        `json:"index_version"`
	Entries      SearchResults `json:"entries"`
}

type RemoteIndexMetadata struct {
	ReleaseURL string `json:"html_url"`
	Version    string `json:"tag_name"`
	Name       string `json:"name"`
}

type Track struct {
	CDNumber     string
	Number       string
	Name         string
	Duration     string
	MP3FileSize  string
	FlacFileSize string
	URL          string
}
