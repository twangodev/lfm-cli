package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/synctest"
	"time"

	"github.com/sony/gobreaker/v2"
	lfm "github.com/twangodev/lfm-api"
)

func TestFallbackCachingAndExpiry(t *testing.T) {
	now := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	observed := now
	status := "playing"
	code := 200
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if code != 200 {
			w.Header().Set("Retry-After", "90")
			w.WriteHeader(code)
			return
		}
		var track any
		if status == "playing" {
			track = map[string]any{"name": "Song", "artist": "Artist", "album": "Album", "url": "https://www.last.fm/music/Artist/_/Song", "image": map[string]string{"extralarge": "https://example.com/cover"}}
		}
		json.NewEncoder(w).Encode(map[string]any{"status": status, "user": "test", "observed_at": observed, "track": track})
	}))
	defer server.Close()
	s := newScrobbleSource(server.URL)
	s.now = func() time.Time { return now }
	s.fallback.now = s.now
	primaryCalls := 0
	s.primary = func(string) (lfm.Scrobble, error) {
		primaryCalls++
		return lfm.EmptyScrobble, fmt.Errorf("status 600")
	}
	first, err := s.get("test")
	if err != nil || !first.Active || first.Name != "Song" || !first.DataTimestamp.Equal(observed) {
		t.Fatalf("%+v %v", first, err)
	}
	now = now.Add(10 * time.Second)
	second, err := s.get("test")
	if err != nil || second.DataTimestamp != first.DataTimestamp || calls != 1 || primaryCalls != 2 {
		t.Fatalf("cache/timestamp failed: %+v %v calls=%d/%d", second, err, calls, primaryCalls)
	}
	now = now.Add(20 * time.Second)
	code = 429
	if _, err = s.get("test"); err != nil {
		t.Fatalf("recent cached result should survive 429: %v", err)
	}
	now = now.Add(60 * time.Second)
	if _, err = s.get("test"); err != nil || calls != 2 || primaryCalls != 2 {
		t.Fatalf("cooldown failed: %v calls=%d/%d", err, calls, primaryCalls)
	}
	now = now.Add(31 * time.Second)
	if _, err = s.get("test"); !errors.Is(err, errScrobbleExpired) {
		t.Fatalf("must expire cached presence: %v", err)
	}
	// Once upstream explicitly reports idle, the cached playing track is cleared.
	now = now.Add(91 * time.Second)
	code = 200
	observed = now
	status = "idle"
	if got, err := s.get("test"); err != nil || got.Active {
		t.Fatalf("idle %+v %v", got, err)
	}

}

func TestFallbackRejectsInvalidData(t *testing.T) {
	now := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	for _, body := range []string{
		`{"status":"idle","user":"test","track":null}`,
		`{"status":"playing","user":"test","observed_at":"2026-09-28T12:00:00Z","track":null}`,
		`{"status":"idle","user":"another","observed_at":"2026-09-28T12:00:00Z","track":null}`,
		`{"status":"idle","user":"test","observed_at":"2026-09-28T11:57:00Z","track":null}`,
		`{"status":"idle","user":"test","observed_at":"2026-09-28T13:00:00Z","track":null}`,
		`<html>Unavailable</html>`,
	} {
		t.Run(body, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, body) }))
			defer server.Close()
			s := newScrobbleSource(server.URL)
			s.now = func() time.Time { return now }
			s.fallback.now = s.now
			if _, err := s.fallback.fetch("test"); err == nil {
				t.Fatal("invalid response accepted")
			}
		})
	}
}

func TestHealthyIdleDoesNotUseFallback(t *testing.T) {
	s := newScrobbleSource("http://unused.invalid")
	s.primary = func(string) (lfm.Scrobble, error) { return lfm.EmptyScrobble, nil }
	if got, err := s.get("test"); err != nil || got.Active || !s.fallback.nextRequest.IsZero() {
		t.Fatalf("%+v %v", got, err)
	}
}

func TestFallbackDisabled(t *testing.T) {
	s := newScrobbleSource("")
	want := errors.New("primary unavailable")
	s.primary = func(string) (lfm.Scrobble, error) { return lfm.EmptyScrobble, want }
	if _, err := s.get("test"); !errors.Is(err, want) {
		t.Fatal(err)
	}
}

func TestWebsiteCircuitBreakerRecovery(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		source := newScrobbleSource("")
		calls := 0
		healthy := false
		source.primary = func(string) (lfm.Scrobble, error) {
			calls++
			if healthy {
				return lfm.EmptyScrobble, nil
			}
			return lfm.EmptyScrobble, errors.New("unavailable")
		}
		source.get("test")
		source.get("test")
		if _, err := source.get("test"); !errors.Is(err, gobreaker.ErrOpenState) || calls != 2 {
			t.Fatalf("breaker did not stop repeated requests: calls=%d err=%v", calls, err)
		}
		time.Sleep(time.Minute + time.Nanosecond)
		healthy = true
		if _, err := source.get("test"); err != nil || calls != 3 {
			t.Fatalf("recovery failed: calls=%d err=%v", calls, err)
		}
		source.get("test")
		if calls != 4 {
			t.Fatal("healthy polling did not resume")
		}
	})
}

func TestWebsiteRecoveryPreservesFallbackRetryAfter(t *testing.T) {
	now := time.Now()
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Retry-After", "90")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()
	source := newScrobbleSource(server.URL)
	source.now = func() time.Time { return now }
	source.fallback.now = source.now
	healthy := false
	source.primary = func(string) (lfm.Scrobble, error) {
		if healthy {
			return lfm.EmptyScrobble, nil
		}
		return lfm.EmptyScrobble, errors.New("website unavailable")
	}
	source.get("test")
	healthy = true
	if _, err := source.get("test"); err != nil {
		t.Fatal(err)
	}
	healthy = false
	now = now.Add(20 * time.Second)
	if _, err := source.get("test"); err == nil {
		t.Fatal("no data during cooldown must be an error")
	}
	if calls != 1 {
		t.Fatal("website recovery discarded fallback Retry-After")
	}
	now = now.Add(71 * time.Second)
	source.get("test")
	if calls != 2 {
		t.Fatal("fallback did not resume after cooldown")
	}
}

func TestPrimaryObservationExpiresWithoutFallback(t *testing.T) {
	now := time.Now()
	source := newScrobbleSource("")
	source.now = func() time.Time { return now }
	healthy := true
	source.primary = func(string) (lfm.Scrobble, error) {
		if healthy {
			return lfm.Scrobble{Active: true, Name: "Song", Artist: "Artist"}, nil
		}
		return lfm.EmptyScrobble, errors.New("website unavailable")
	}
	if _, err := source.get("test"); err != nil {
		t.Fatal(err)
	}
	healthy = false
	now = now.Add(maxScrobbleAge - time.Second)
	if _, err := source.get("test"); errors.Is(err, errScrobbleExpired) {
		t.Fatal("fresh presence expired too soon")
	}
	now = now.Add(time.Second)
	if _, err := source.get("test"); !errors.Is(err, errScrobbleExpired) {
		t.Fatalf("primary observation did not expire: %v", err)
	}
}

func TestSourceChangesUseObservationAge(t *testing.T) {
	now := time.Now()
	observed := now.Add(-90 * time.Second)
	available := true
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !available {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"status": "playing", "user": "test", "observed_at": observed,
			"track": map[string]any{"name": "Song", "artist": "Artist", "album": "Album"},
		})
	}))
	defer server.Close()
	source := newScrobbleSource(server.URL)
	source.now = func() time.Time { return now }
	source.fallback.now = source.now
	healthy := false
	source.primary = func(string) (lfm.Scrobble, error) {
		if healthy {
			return lfm.Scrobble{Active: true, Name: "Song", Artist: "Artist", Album: "Album"}, nil
		}
		return lfm.EmptyScrobble, errors.New("website unavailable")
	}
	first, err := source.get("test")
	if err != nil || !first.DataTimestamp.Equal(observed) {
		t.Fatalf("fallback result: %+v %v", first, err)
	}
	healthy = true
	now = now.Add(10 * time.Second)
	if _, err = source.get("test"); err != nil {
		t.Fatal(err)
	}
	available = false
	healthy = false
	now = now.Add(30 * time.Second)
	if _, err = source.get("test"); err == nil || errors.Is(err, errScrobbleExpired) {
		t.Fatalf("recent primary result must survive fallback expiry: %v", err)
	}
	now = now.Add(90 * time.Second)
	if _, err = source.get("test"); !errors.Is(err, errScrobbleExpired) {
		t.Fatalf("source observation did not expire: %v", err)
	}
	// The same song after expiry is a new observation, not a continued timer.
	available = true
	observed = now
	now = now.Add(fallbackPollInterval)
	track, err := source.get("test")
	if err != nil || track.DataTimestamp.Equal(first.DataTimestamp) {
		t.Fatalf("expired elapsed timer was reused: %+v %v", track, err)
	}
}
