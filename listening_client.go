package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	lfm "github.com/twangodev/lfm-api"
)

const fallbackPollInterval = 30 * time.Second
const maxListeningResponseBytes = 64 * 1024

type listeningClient struct {
	baseURL     string
	http        *http.Client
	now         func() time.Time
	nextRequest time.Time
	cached      scrobbleObservation
	lastError   error
}

func newListeningClient(endpoint string) *listeningClient {
	return &listeningClient{baseURL: strings.TrimRight(endpoint, "/"), http: &http.Client{Timeout: 10 * time.Second}, now: time.Now}
}

func (c *listeningClient) reset() {
	c.nextRequest = time.Time{}
	c.invalidateTrack()
	c.lastError = nil
}

func (c *listeningClient) invalidateTrack() {
	c.cached = scrobbleObservation{}
}

func (c *listeningClient) get(username string) (scrobbleObservation, error) {
	if !c.now().Before(c.nextRequest) {
		observation, err := c.fetch(username)
		c.lastError = err
		if err == nil {
			c.cached = observation
		}
	}
	if c.cached.fresh(c.now()) {
		return c.cached, nil
	}
	if c.lastError == nil {
		return scrobbleObservation{}, errors.New("waiting for listening service retry window")
	}
	return scrobbleObservation{}, c.lastError
}

type playingResponse struct {
	Status     string        `json:"status"`
	User       string        `json:"user"`
	ObservedAt time.Time     `json:"observed_at"`
	Track      *playingTrack `json:"track"`
}

type playingTrack struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	URL    string `json:"url"`
	Image  struct {
		ExtraLarge string `json:"extralarge"`
		Large      string `json:"large"`
	} `json:"image"`
}

func (c *listeningClient) fetch(username string) (scrobbleObservation, error) {
	now := c.now()
	c.nextRequest = now.Add(fallbackPollInterval)
	response, err := c.http.Get(c.baseURL + "/playing/" + url.PathEscape(username))
	if err != nil {
		return scrobbleObservation{}, err
	}
	defer response.Body.Close()
	if until := retryAfterTime(response.Header.Get("Retry-After"), now); until.After(c.nextRequest) {
		c.nextRequest = until
	}
	if response.StatusCode != http.StatusOK {
		return scrobbleObservation{}, fmt.Errorf("listening service returned status %d", response.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, maxListeningResponseBytes+1))
	if err != nil {
		return scrobbleObservation{}, err
	}
	if len(data) > maxListeningResponseBytes {
		return scrobbleObservation{}, errors.New("listening response exceeds size limit")
	}
	var body playingResponse
	if err = json.Unmarshal(data, &body); err != nil {
		return scrobbleObservation{}, fmt.Errorf("decode listening response: %w", err)
	}
	return body.observation(username, c.now())
}

func retryAfterTime(header string, now time.Time) time.Time {
	if seconds, err := strconv.ParseUint(header, 10, 32); err == nil {
		return now.Add(time.Duration(seconds) * time.Second)
	}
	until, _ := http.ParseTime(header)
	return until
}

func (body playingResponse) observation(username string, now time.Time) (scrobbleObservation, error) {
	result := scrobbleObservation{observedAt: body.ObservedAt}
	if !strings.EqualFold(body.User, username) || !result.fresh(now) || body.ObservedAt.After(now.Add(30*time.Second)) {
		return scrobbleObservation{}, errors.New("listening response has invalid or expired freshness data")
	}
	if body.Status == "idle" && body.Track == nil {
		return result, nil
	}
	if body.Status != "playing" || body.Track == nil || body.Track.Name == "" || body.Track.Artist == "" {
		return scrobbleObservation{}, errors.New("invalid listening track response")
	}
	track := body.Track
	cover := track.Image.ExtraLarge
	if cover == "" {
		cover = track.Image.Large
	}
	result.track = lfm.Scrobble{Active: true, Name: track.Name, Artist: track.Artist, Album: track.Album, DataLink: track.URL, DataLinkTitle: "View track on Last.fm", CoverArtUrl: cover}
	return result, nil
}
