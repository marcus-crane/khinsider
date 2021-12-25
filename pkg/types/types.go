package types

type Album struct {
  FlacAvailable bool
  Name string
  MP3FileSize string
  FlacFileSize string
  FileCount int
  Tracks []Track
}

type SearchResults map[string]string

type Track struct {
  SongName string
  SongURL  string
}
