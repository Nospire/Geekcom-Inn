package jukebox

import (
	_ "embed"
	"math/rand/v2"
	"net/url"
	"strings"
)

//go:embed jazz_catalog.txt
var jazzCatalogRaw string

//go:embed electronic_catalog.txt
var electronicCatalogRaw string

//go:embed cantina_catalog.txt
var cantinaCatalogRaw string

const archiveDownloadBase = "https://archive.org/download/"

// Catalog holds all genre track lists.
type Catalog struct {
	lofi       *Lofi
	jazz       []genreTrack
	electronic []genreTrack
	cantina    []genreTrack
}

// genreTrack is a track from an archive.org collection.
type genreTrack struct {
	path  string // collection-id/filename (raw, with spaces)
	title string
}

func NewCatalog() *Catalog {
	return &Catalog{
		lofi:       NewLofi(),
		jazz:       parseArchiveGenreCatalog(jazzCatalogRaw),
		electronic: parseArchiveGenreCatalog(electronicCatalogRaw),
		cantina:    parseArchiveGenreCatalog(cantinaCatalogRaw),
	}
}

// TrackCount returns the number of tracks for a genre.
func (c *Catalog) TrackCount(g Genre) int {
	switch g {
	case GenreLofi:
		return c.lofi.TrackCount()
	case GenreJazz:
		return len(c.jazz)
	case GenreElectronic:
		return len(c.electronic)
	case GenreCantina:
		return len(c.cantina)
	}
	return 0
}

// RandomTracks picks n random tracks from the given genre.
func (c *Catalog) RandomTracks(g Genre, n int) []Track {
	switch g {
	case GenreLofi:
		return c.lofi.randomTracks(n)
	case GenreJazz:
		return pickArchiveTracks(c.jazz, n, "Jazz")
	case GenreElectronic:
		return pickArchiveTracks(c.electronic, n, "Electronic")
	case GenreCantina:
		return pickArchiveTracks(c.cantina, n, "Cantina")
	}
	return nil
}

func pickArchiveTracks(catalog []genreTrack, n int, artist string) []Track {
	if n > len(catalog) {
		n = len(catalog)
	}
	if n == 0 {
		return nil
	}
	indices := make([]int, len(catalog))
	for i := range indices {
		indices[i] = i
	}
	for i := 0; i < n; i++ {
		j := i + rand.IntN(len(indices)-i)
		indices[i], indices[j] = indices[j], indices[i]
	}
	tracks := make([]Track, n)
	for i := 0; i < n; i++ {
		t := catalog[indices[i]]
		tracks[i] = Track{
			ID:     t.path,
			Title:  t.title,
			Artist: artist,
			URL:    archiveDownloadBase + encodeArchivePath(t.path),
		}
	}
	return tracks
}

// encodeArchivePath URL-encodes each segment of a catalog path.
// Input: "JazzClassics/001  Idris Muhammad - Power Of Soul.mp3"
// Output: "JazzClassics/001%20%20Idris%20Muhammad%20-%20Power%20Of%20Soul.mp3"
func encodeArchivePath(raw string) string {
	parts := strings.SplitN(raw, "/", 2)
	if len(parts) != 2 {
		return url.PathEscape(raw)
	}
	return parts[0] + "/" + url.PathEscape(parts[1])
}

// parseArchiveGenreCatalog parses lines of "collection-id/filename.mp3".
func parseArchiveGenreCatalog(raw string) []genreTrack {
	var tracks []genreTrack
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasSuffix(strings.ToLower(line), ".mp3") {
			continue
		}
		// Extract filename without extension for title
		parts := strings.Split(line, "/")
		name := parts[len(parts)-1]
		name = strings.TrimSuffix(name, ".mp3")
		// Remove track number prefix like "01 " or "01. "
		if len(name) > 3 && name[0] >= '0' && name[0] <= '9' {
			for i, c := range name {
				if c == ' ' || c == '-' {
					name = strings.TrimSpace(name[i+1:])
					break
				}
				if i > 3 {
					break
				}
			}
		}
		tracks = append(tracks, genreTrack{
			path:  line,
			title: name,
		})
	}
	return tracks
}
