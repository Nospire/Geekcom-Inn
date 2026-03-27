package jukebox

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJamendoSearch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3.0/tracks/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("search") != "lofi" {
			t.Errorf("expected search=lofi, got %s", q.Get("search"))
		}
		if q.Get("client_id") != "test-id" {
			t.Errorf("expected client_id=test-id, got %s", q.Get("client_id"))
		}
		resp := jamendoResponse{
			Results: []jamendoTrack{
				{
					ID:       "12345",
					Name:     "Chill Sunset",
					Artist:   "ambient_collective",
					Duration: 210,
					Audio:    "https://mp3.jamendo.com/track/12345.mp3",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	j := &Jamendo{clientID: "test-id", baseURL: ts.URL}
	tracks, err := j.Search(context.Background(), "lofi", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("expected 1 track, got %d", len(tracks))
	}
	if tracks[0].Title != "Chill Sunset" {
		t.Errorf("expected title 'Chill Sunset', got '%s'", tracks[0].Title)
	}
	if tracks[0].Source != "jamendo" {
		t.Errorf("expected source 'jamendo', got '%s'", tracks[0].Source)
	}
}

func TestJamendoSearchEmpty(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := jamendoResponse{Results: []jamendoTrack{}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	j := &Jamendo{clientID: "test-id", baseURL: ts.URL}
	tracks, err := j.Search(context.Background(), "nonexistent", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tracks) != 0 {
		t.Errorf("expected 0 tracks, got %d", len(tracks))
	}
}

func TestJamendoStreamURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := jamendoResponse{
			Results: []jamendoTrack{
				{
					ID:    "12345",
					Audio: "https://mp3.jamendo.com/track/12345.mp3",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	j := &Jamendo{clientID: "test-id", baseURL: ts.URL}
	url, err := j.StreamURL(context.Background(), "12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url != "https://mp3.jamendo.com/track/12345.mp3" {
		t.Errorf("unexpected URL: %s", url)
	}
}
