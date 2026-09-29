package main

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/sony/gobreaker/v2"
	lfm "github.com/twangodev/lfm-api"
)

const maxScrobbleAge = 2 * time.Minute

var errScrobbleExpired = errors.New("listening data expired")

type scrobbleObservation struct {
	track      lfm.Scrobble
	observedAt time.Time
}

func (o scrobbleObservation) fresh(now time.Time) bool {
	return !o.observedAt.IsZero() && now.Sub(o.observedAt) < maxScrobbleAge
}

type scrobbleSource struct {
	primary       func(string) (lfm.Scrobble, error)
	breaker       *gobreaker.CircuitBreaker[lfm.Scrobble]
	fallback      *listeningClient
	last          scrobbleObservation
	now           func() time.Time
	username      string
	usingFallback bool
}

func newWebsiteBreaker() *gobreaker.CircuitBreaker[lfm.Scrobble] {
	return gobreaker.NewCircuitBreaker[lfm.Scrobble](gobreaker.Settings{
		Name:        "lastfm-website",
		Timeout:     time.Minute,
		ReadyToTrip: func(counts gobreaker.Counts) bool { return counts.ConsecutiveFailures >= 2 },
		OnStateChange: func(name string, from, to gobreaker.State) {
			slog.Debug("Website circuit breaker changed state.", "from", from.String(), "to", to.String())
		},
	})
}

func newScrobbleSource(endpoint string) *scrobbleSource {
	return &scrobbleSource{
		primary:  lfm.GetActiveScrobble,
		breaker:  newWebsiteBreaker(),
		fallback: newListeningClient(endpoint),
		now:      time.Now,
	}
}

func (s *scrobbleSource) get(username string) (lfm.Scrobble, error) {
	if s.username != username {
		s.username = username
		s.last = scrobbleObservation{}
		s.usingFallback = false
		s.breaker = newWebsiteBreaker()
		s.fallback.reset()
	}
	track, websiteErr := s.breaker.Execute(func() (lfm.Scrobble, error) { return s.primary(username) })
	if websiteErr == nil {
		s.fallback.invalidateTrack()
		if s.usingFallback {
			slog.Info("Last.fm website fetching recovered.")
			s.usingFallback = false
		}
		return s.remember(scrobbleObservation{track: track, observedAt: s.now()}), nil
	}
	slog.Debug("Last.fm website fetch unavailable.", "error", websiteErr)
	if s.fallback.baseURL == "" {
		return s.unavailable(websiteErr)
	}
	observation, fallbackErr := s.fallback.get(username)
	if fallbackErr != nil {
		return s.unavailable(fmt.Errorf("website: %v; fallback: %w", websiteErr, fallbackErr))
	}
	if !s.usingFallback {
		slog.Info("Using listening fallback.", "endpoint", s.fallback.baseURL)
		s.usingFallback = true
	}
	return s.remember(observation), nil
}

func (s *scrobbleSource) unavailable(err error) (lfm.Scrobble, error) {
	if !s.last.observedAt.IsZero() && !s.last.fresh(s.now()) {
		return lfm.EmptyScrobble, fmt.Errorf("%w: %w", errScrobbleExpired, err)
	}
	return lfm.EmptyScrobble, err
}

func (s *scrobbleSource) remember(observation scrobbleObservation) lfm.Scrobble {
	track := observation.track
	// Now-playing API responses have no start time. Preserve the first observation
	// while a fresh track stays the same so cache hits do not reset Discord's timer.
	if track.Active && track.DataTimestamp.IsZero() {
		previous := s.last.track
		if s.last.fresh(s.now()) && previous.Active && track.Name == previous.Name && track.Artist == previous.Artist && track.Album == previous.Album {
			track.DataTimestamp = previous.DataTimestamp
		} else {
			track.DataTimestamp = observation.observedAt
		}
	}
	observation.track = track
	s.last = observation
	return track
}
