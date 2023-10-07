package types

type SearchResults map[string]string

type Album struct {
	Title             string            `json:"title"`
	TitlesAlternative []string          `json:"titles_alternative"`
	CrawledAt         string            `json:"crawled_at"`
	Slug              string            `json:"slug"`
	URL               string            `json:"url"`
	Platforms         map[string]string `json:"platforms"`
	Images            []string          `json:"images"`
	Year              int32             `json:"year"`
	DateAdded         string            `json:"date_added"`
	Developers        map[string]string `json:"developers"`
	Publishers        map[string]string `json:"publishers"`
	Genres            map[string]string `json:"genres"`
	Tracks            []Track           `json:"tracks"`
	Total             Total             `json:"total"`
}

type Track struct {
	DiscNumber        int32  `json:"disc_number"`
	FilesizeMP3Bytes  int64  `json:"filesize_mp3_bytes"`
	FilesizeFlacBytes int64  `json:"filesize_flac_bytes"`
	TrackNumber       int32  `json:"track_number"`
	Title             string `json:"title"`
	Runtime           int32  `json:"runtime"`
	SourceMP3         string `json:"source_mp3"`
	SourceFlac        string `json:"source_flac"`
	TrackURL          string `json:"track_url"`
}

type Total struct {
	Runtime           int32 `json:"runtime"`
	FilesizeMP3Bytes  int64 `json:"filesize_mp3_bytes"`
	FilesizeFlacBytes int64 `json:"filesize_flac_bytes"`
	Tracks            int32 `json:"tracks"`
}

type RemoteIndexMetadata struct {
	ReleaseURL string `json:"html_url"`
	Version    string `json:"tag_name"`
	Name       string `json:"name"`
	Prerelease bool   `json:"prerelease"`
}
