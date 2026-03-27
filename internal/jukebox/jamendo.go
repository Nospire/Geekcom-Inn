package jukebox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const jamendoBaseURL = "https://api.jamendo.com"

type jamendoTrack struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist_name"`
	Duration int    `json:"duration"`
	Audio    string `json:"audio"`
}

type jamendoResponse struct {
	Results []jamendoTrack `json:"results"`
}

// Jamendo implements MusicBackend for the Jamendo API.
type Jamendo struct {
	clientID string
	baseURL  string
	client   *http.Client
}

// NewJamendo creates a Jamendo backend. clientID is from devportal.jamendo.com.
func NewJamendo(clientID string) *Jamendo {
	return &Jamendo{
		clientID: clientID,
		baseURL:  jamendoBaseURL,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (j *Jamendo) Name() string { return "jamendo" }

func (j *Jamendo) Enabled() bool { return j.clientID != "" }

func (j *Jamendo) httpClient() *http.Client {
	if j.client != nil {
		return j.client
	}
	return http.DefaultClient
}

func (j *Jamendo) Search(ctx context.Context, query string, limit int) ([]Track, error) {
	u, _ := url.Parse(j.baseURL + "/v3.0/tracks/")
	q := u.Query()
	q.Set("client_id", j.clientID)
	q.Set("search", query)
	q.Set("limit", strconv.Itoa(limit))
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("jamendo: build request: %w", err)
	}

	resp, err := j.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("jamendo: search: %w", err)
	}
	defer resp.Body.Close()

	var jr jamendoResponse
	if err := json.NewDecoder(resp.Body).Decode(&jr); err != nil {
		return nil, fmt.Errorf("jamendo: decode: %w", err)
	}

	tracks := make([]Track, 0, len(jr.Results))
	for _, t := range jr.Results {
		tracks = append(tracks, Track{
			ID:       t.ID,
			Title:    t.Name,
			Artist:   t.Artist,
			Duration: t.Duration,
			URL:      t.Audio,
			Source:   "jamendo",
		})
	}
	return tracks, nil
}

func (j *Jamendo) StreamURL(ctx context.Context, trackID string) (string, error) {
	u, _ := url.Parse(j.baseURL + "/v3.0/tracks/")
	q := u.Query()
	q.Set("client_id", j.clientID)
	q.Set("id", trackID)
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("jamendo: build request: %w", err)
	}

	resp, err := j.httpClient().Do(req)
	if err != nil {
		return "", fmt.Errorf("jamendo: stream url: %w", err)
	}
	defer resp.Body.Close()

	var jr jamendoResponse
	if err := json.NewDecoder(resp.Body).Decode(&jr); err != nil {
		return "", fmt.Errorf("jamendo: decode: %w", err)
	}
	if len(jr.Results) == 0 {
		return "", fmt.Errorf("jamendo: track %s not found", trackID)
	}
	return jr.Results[0].Audio, nil
}
