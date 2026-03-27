package jukebox

import (
	"testing"
	"time"
)

func TestEngineInitialState(t *testing.T) {
	e := NewEngine(nil)
	state := e.State()
	if state.Phase != PhaseIdle {
		t.Errorf("expected PhaseIdle, got %v", state.Phase)
	}
	if state.Current != nil {
		t.Errorf("expected no current track")
	}
}

func TestEngineAddRequest(t *testing.T) {
	e := NewEngine(nil)
	e.mu.Lock()
	e.phase = PhaseRequesting
	e.mu.Unlock()

	track := Track{ID: "1", Title: "Test", Artist: "Artist", Source: "jamendo"}
	e.AddRequest("user1", track)

	state := e.State()
	if len(state.Requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(state.Requests))
	}
	if state.Requests[0].Count != 1 {
		t.Errorf("expected count 1, got %d", state.Requests[0].Count)
	}
}

func TestEngineAddRequestDuplicate(t *testing.T) {
	e := NewEngine(nil)
	e.mu.Lock()
	e.phase = PhaseRequesting
	e.mu.Unlock()

	track := Track{ID: "1", Title: "Test", Artist: "Artist", Source: "jamendo"}
	e.AddRequest("user1", track)
	e.AddRequest("user2", track)

	state := e.State()
	if len(state.Requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(state.Requests))
	}
	if state.Requests[0].Count != 2 {
		t.Errorf("expected count 2, got %d", state.Requests[0].Count)
	}
}

func TestEngineVote(t *testing.T) {
	e := NewEngine(nil)
	track := Track{ID: "1", Title: "Test", Artist: "Artist", Source: "jamendo"}

	e.mu.Lock()
	e.phase = PhaseVoting
	e.shortlist = []Track{track}
	e.votes = make(map[string]map[string]bool)
	e.mu.Unlock()

	ok := e.Vote("user1", "1")
	if !ok {
		t.Error("expected vote to succeed")
	}

	ok = e.Vote("user1", "1")
	if ok {
		t.Error("expected duplicate vote to fail")
	}
}

func TestEngineVoteNotInShortlist(t *testing.T) {
	e := NewEngine(nil)
	track := Track{ID: "1", Title: "Test", Artist: "Artist", Source: "jamendo"}

	e.mu.Lock()
	e.phase = PhaseVoting
	e.shortlist = []Track{track}
	e.votes = make(map[string]map[string]bool)
	e.mu.Unlock()

	ok := e.Vote("user1", "nonexistent")
	if ok {
		t.Error("expected vote for nonexistent track to fail")
	}
}

func TestEngineBuildShortlist(t *testing.T) {
	e := NewEngine(nil)

	e.mu.Lock()
	e.phase = PhaseRequesting
	e.requestPool = map[string]*Request{
		"a": {Track: Track{ID: "a", Title: "A"}, Count: 5},
		"b": {Track: Track{ID: "b", Title: "B"}, Count: 3},
		"c": {Track: Track{ID: "c", Title: "C"}, Count: 7},
		"d": {Track: Track{ID: "d", Title: "D"}, Count: 1},
		"e": {Track: Track{ID: "e", Title: "E"}, Count: 4},
		"f": {Track: Track{ID: "f", Title: "F"}, Count: 2},
		"g": {Track: Track{ID: "g", Title: "G"}, Count: 6},
	}
	e.mu.Unlock()

	shortlist := e.buildShortlist()
	if len(shortlist) != 5 {
		t.Fatalf("expected 5 tracks in shortlist, got %d", len(shortlist))
	}
	if shortlist[0].ID != "c" {
		t.Errorf("expected first track to be 'c', got '%s'", shortlist[0].ID)
	}
	if shortlist[4].ID != "b" {
		t.Errorf("expected last track to be 'b', got '%s'", shortlist[4].ID)
	}
}

func TestEnginePickWinner(t *testing.T) {
	e := NewEngine(nil)

	tracks := []Track{
		{ID: "a", Title: "A"},
		{ID: "b", Title: "B"},
	}
	e.mu.Lock()
	e.phase = PhaseVoting
	e.shortlist = tracks
	e.votes = map[string]map[string]bool{
		"a": {"user1": true},
		"b": {"user2": true, "user3": true},
	}
	e.mu.Unlock()

	winner := e.pickWinner()
	if winner == nil {
		t.Fatal("expected a winner")
	}
	if winner.ID != "b" {
		t.Errorf("expected winner 'b' (2 votes), got '%s'", winner.ID)
	}
}

func TestEngineTickPhaseTransitions(t *testing.T) {
	e := NewEngine(nil)

	track := Track{ID: "1", Title: "Test", Duration: 100}
	e.mu.Lock()
	e.current = &track
	e.playStart = time.Now().Add(-76 * time.Second)
	e.phase = PhasePlaying
	e.mu.Unlock()

	e.tick()

	state := e.State()
	if state.Phase != PhaseRequesting {
		t.Errorf("expected PhaseRequesting at 76%%, got %v", state.Phase)
	}

	e.mu.Lock()
	e.playStart = time.Now().Add(-91 * time.Second)
	e.mu.Unlock()

	e.tick()

	state = e.State()
	if state.Phase != PhaseVoting {
		t.Errorf("expected PhaseVoting at 91%%, got %v", state.Phase)
	}

	e.mu.Lock()
	e.playStart = time.Now().Add(-101 * time.Second)
	e.mu.Unlock()

	e.tick()

	state = e.State()
	if state.Phase != PhaseIdle {
		t.Errorf("expected PhaseIdle after track ends with no votes/backends, got %v", state.Phase)
	}
}

func TestEngineFinishTrack(t *testing.T) {
	e := NewEngine(nil)
	track := Track{ID: "1", Title: "Test"}
	nextTrack := Track{ID: "2", Title: "Next"}

	e.mu.Lock()
	e.current = &track
	e.phase = PhaseVoting
	e.shortlist = []Track{nextTrack}
	e.votes = map[string]map[string]bool{
		"2": {"user1": true},
	}
	e.mu.Unlock()

	winner := e.FinishTrack()
	if winner == nil || winner.ID != "2" {
		t.Errorf("expected winner '2', got %v", winner)
	}

	state := e.State()
	if state.Phase != PhasePlaying {
		t.Errorf("expected PhasePlaying after finish, got %v", state.Phase)
	}
	if state.Current == nil || state.Current.ID != "2" {
		t.Error("expected current track to be the winner")
	}
}
