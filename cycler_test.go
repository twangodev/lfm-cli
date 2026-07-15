package main

import (
	"errors"
	"testing"
	"time"

	lfm "github.com/twangodev/lfm-api"
)

func TestCycleRetainsLoginOnFetchError(t *testing.T) {
	original := getActiveScrobble
	getActiveScrobble = func(string) (lfm.Scrobble, error) {
		return lfm.EmptyScrobble, errors.New("last.fm returned status 600")
	}
	defer func() {
		getActiveScrobble = original
		loggedIn = false
	}()

	loggedIn = true
	cycle()

	if !loggedIn {
		t.Fatal("transient fetch error must not log out of Discord RPC")
	}
}

func TestLogoutResetsTimestamp(t *testing.T) {
	ts = time.Now()
	defer func() { loggedIn = false }()

	logout()

	if !ts.IsZero() {
		t.Fatal("logout must reset ts so re-login re-sends presence for the same track")
	}
}
