package main

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	lfm "github.com/twangodev/lfm-api"
)

func TestCycleRetainsLoginOnFetchError(t *testing.T) {
	original := getActiveScrobble
	originalHealth := fetchHealth
	getActiveScrobble = func(string) (lfm.Scrobble, error) {
		return lfm.EmptyScrobble, errors.New("last.fm returned status 600")
	}
	defer func() {
		getActiveScrobble = original
		fetchHealth = originalHealth
		loggedIn = false
	}()

	loggedIn = true
	cycle()

	if !loggedIn {
		t.Fatal("transient fetch error must not log out of Discord RPC")
	}
}

func TestScrobbleFetchHealthLogging(t *testing.T) {
	var logs bytes.Buffer
	originalLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(originalLogger) })

	var health scrobbleFetchHealth
	start := time.Now()
	fetchErr := errors.New("last.fm returned status 600")
	observe := func(seconds int, err error) {
		health.observe(start.Add(time.Duration(seconds)*time.Second), err)
	}
	assertCounts := func(debug, warn, info int) {
		t.Helper()
		for level, want := range map[string]int{"DEBUG": debug, "WARN": warn, "INFO": info} {
			if got := strings.Count(logs.String(), "level="+level); got != want {
				t.Fatalf("%s count = %d, want %d; logs:\n%s", level, got, want, &logs)
			}
		}
	}

	observe(0, fetchErr)
	observe(59, fetchErr)
	assertCounts(2, 0, 0)
	observe(60, nil) // A brief failure recovers silently and resets the timer.
	assertCounts(2, 0, 0)
	observe(61, fetchErr)
	observe(120, fetchErr)
	assertCounts(4, 0, 0)
	observe(121, fetchErr)
	observe(181, fetchErr) // Continued failure must not repeat the warning.
	assertCounts(6, 1, 0)
	observe(182, nil)
	observe(183, nil)
	assertCounts(6, 1, 1)
	observe(184, fetchErr) // A new outage gets its own warning and recovery.
	observe(244, fetchErr)
	observe(245, nil)
	assertCounts(8, 2, 2)
}

func TestLogoutResetsTimestamp(t *testing.T) {
	ts = time.Now()
	defer func() { loggedIn = false }()

	logout()

	if !ts.IsZero() {
		t.Fatal("logout must reset ts so re-login re-sends presence for the same track")
	}
}
