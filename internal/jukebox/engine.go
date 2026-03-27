package jukebox

import (
	"context"
	"math/rand/v2"
	"sort"
	"sync"
	"time"
)

type Phase int

const (
	PhaseIdle       Phase = iota
	PhasePlaying
	PhaseRequesting
	PhaseVoting
)

type Request struct {
	Track     Track
	Count     int
	FirstTime time.Time
}

type VoteTally struct {
	Track Track
	Votes int
}

type EngineState struct {
	Phase     Phase
	Current   *Track
	Position  time.Duration
	Requests  []Request
	Shortlist []VoteTally
	Listeners int
}

type OnStateChange func()

type Engine struct {
	mu            sync.RWMutex
	phase         Phase
	current       *Track
	playStart     time.Time
	requestPool   map[string]*Request
	shortlist     []Track
	votes         map[string]map[string]bool
	userVoted     map[string]bool
	backends      []MusicBackend
	onStateChange OnStateChange
	listeners     int
}

func NewEngine(backends []MusicBackend) *Engine {
	return &Engine{
		phase:       PhaseIdle,
		requestPool: make(map[string]*Request),
		votes:       make(map[string]map[string]bool),
		userVoted:   make(map[string]bool),
		backends:    backends,
	}
}

func (e *Engine) SetOnStateChange(fn OnStateChange) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.onStateChange = fn
}

func (e *Engine) SetListeners(n int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners = n
}

func (e *Engine) Backends() []MusicBackend {
	var enabled []MusicBackend
	for _, b := range e.backends {
		if b.Enabled() {
			enabled = append(enabled, b)
		}
	}
	return enabled
}

func (e *Engine) State() EngineState {
	e.mu.RLock()
	defer e.mu.RUnlock()

	state := EngineState{
		Phase:     e.phase,
		Current:   e.current,
		Listeners: e.listeners,
	}

	if e.current != nil {
		state.Position = time.Since(e.playStart)
	}

	reqs := make([]Request, 0, len(e.requestPool))
	for _, r := range e.requestPool {
		reqs = append(reqs, *r)
	}
	sort.Slice(reqs, func(i, j int) bool {
		if reqs[i].Count != reqs[j].Count {
			return reqs[i].Count > reqs[j].Count
		}
		return reqs[i].FirstTime.Before(reqs[j].FirstTime)
	})
	state.Requests = reqs

	for _, t := range e.shortlist {
		tally := VoteTally{Track: t}
		if voters, ok := e.votes[t.ID]; ok {
			tally.Votes = len(voters)
		}
		state.Shortlist = append(state.Shortlist, tally)
	}
	sort.Slice(state.Shortlist, func(i, j int) bool {
		return state.Shortlist[i].Votes > state.Shortlist[j].Votes
	})

	return state
}

func (e *Engine) AddRequest(userFP string, track Track) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.phase != PhaseRequesting && e.phase != PhaseVoting {
		return
	}

	if existing, ok := e.requestPool[track.ID]; ok {
		existing.Count++
	} else {
		e.requestPool[track.ID] = &Request{
			Track:     track,
			Count:     1,
			FirstTime: time.Now(),
		}
	}

	e.notifyChange()
}

func (e *Engine) Vote(userFP string, trackID string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.phase != PhaseVoting {
		return false
	}

	if e.userVoted[userFP] {
		return false
	}

	inShortlist := false
	for _, t := range e.shortlist {
		if t.ID == trackID {
			inShortlist = true
			break
		}
	}
	if !inShortlist {
		return false
	}

	if e.votes[trackID] == nil {
		e.votes[trackID] = make(map[string]bool)
	}
	e.votes[trackID][userFP] = true
	e.userVoted[userFP] = true

	e.notifyChange()
	return true
}

func (e *Engine) StartPlaying(track Track) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.current = &track
	e.playStart = time.Now()
	e.phase = PhasePlaying

	e.notifyChange()
}

func (e *Engine) OpenRequesting() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.phase = PhaseRequesting
	e.requestPool = make(map[string]*Request)

	e.notifyChange()
}

func (e *Engine) OpenVoting() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.shortlist = e.buildShortlistLocked()
	e.votes = make(map[string]map[string]bool)
	e.userVoted = make(map[string]bool)
	e.phase = PhaseVoting

	e.notifyChange()
}

func (e *Engine) FinishTrack() *Track {
	e.mu.Lock()
	defer e.mu.Unlock()

	winner := e.pickWinnerLocked()

	e.requestPool = make(map[string]*Request)
	e.shortlist = nil
	e.votes = make(map[string]map[string]bool)
	e.userVoted = make(map[string]bool)

	if winner != nil {
		e.current = winner
		e.playStart = time.Now()
		e.phase = PhasePlaying
	} else {
		e.phase = PhaseIdle
		e.current = nil
	}

	e.notifyChange()
	return winner
}

func (e *Engine) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.tick()
		}
	}
}

func (e *Engine) tick() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.current == nil || e.phase == PhaseIdle {
		e.tryAutoPlay()
		return
	}

	elapsed := time.Since(e.playStart)
	duration := time.Duration(e.current.Duration) * time.Second
	if duration == 0 {
		return
	}

	progress := elapsed.Seconds() / duration.Seconds()

	switch e.phase {
	case PhasePlaying:
		if progress >= 0.75 {
			e.phase = PhaseRequesting
			e.requestPool = make(map[string]*Request)
			e.notifyChange()
		}
	case PhaseRequesting:
		if progress >= 0.90 {
			e.shortlist = e.buildShortlistLocked()
			e.votes = make(map[string]map[string]bool)
			e.userVoted = make(map[string]bool)
			e.phase = PhaseVoting
			e.notifyChange()
		}
	case PhaseVoting:
		if progress >= 1.0 {
			winner := e.pickWinnerLocked()
			e.requestPool = make(map[string]*Request)
			e.shortlist = nil
			e.votes = make(map[string]map[string]bool)
			e.userVoted = make(map[string]bool)

			if winner != nil {
				e.current = winner
				e.playStart = time.Now()
				e.phase = PhasePlaying
			} else {
				e.phase = PhaseIdle
				e.current = nil
				e.tryAutoPlay()
			}
			e.notifyChange()
		}
	}
}

func (e *Engine) tryAutoPlay() {
	for _, b := range e.backends {
		if !b.Enabled() {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		tracks, err := b.Search(ctx, "popular", 20)
		cancel()
		if err != nil || len(tracks) == 0 {
			continue
		}
		pick := tracks[rand.IntN(len(tracks))]
		e.current = &pick
		e.playStart = time.Now()
		e.phase = PhasePlaying
		e.notifyChange()
		return
	}
}

func (e *Engine) buildShortlist() []Track {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.buildShortlistLocked()
}

func (e *Engine) buildShortlistLocked() []Track {
	reqs := make([]*Request, 0, len(e.requestPool))
	for _, r := range e.requestPool {
		reqs = append(reqs, r)
	}
	sort.Slice(reqs, func(i, j int) bool {
		if reqs[i].Count != reqs[j].Count {
			return reqs[i].Count > reqs[j].Count
		}
		return reqs[i].FirstTime.Before(reqs[j].FirstTime)
	})

	limit := 5
	if len(reqs) < limit {
		limit = len(reqs)
	}

	tracks := make([]Track, limit)
	for i := 0; i < limit; i++ {
		tracks[i] = reqs[i].Track
	}
	return tracks
}

func (e *Engine) pickWinner() *Track {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.pickWinnerLocked()
}

func (e *Engine) pickWinnerLocked() *Track {
	if len(e.shortlist) == 0 {
		return nil
	}

	type candidate struct {
		track Track
		votes int
	}
	var candidates []candidate
	for _, t := range e.shortlist {
		v := 0
		if voters, ok := e.votes[t.ID]; ok {
			v = len(voters)
		}
		candidates = append(candidates, candidate{track: t, votes: v})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].votes > candidates[j].votes
	})

	maxVotes := candidates[0].votes
	var tied []candidate
	for _, c := range candidates {
		if c.votes == maxVotes {
			tied = append(tied, c)
		}
	}

	winner := tied[rand.IntN(len(tied))]
	return &winner.track
}

func (e *Engine) notifyChange() {
	if e.onStateChange != nil {
		e.onStateChange()
	}
}
