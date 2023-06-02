package main

import (
	"fmt"
	"github.com/hugolgst/rich-go/client"
	lfm "github.com/twangodev/lfm-api"
	"golang.org/x/net/html"
)

func createActivity(s lfm.Scrobble, songLink bool) client.Activity {

	var songButton *client.Button
	// Determines whether to display profile link in buttons
	if songLink {
		dataLinkTitle := s.DataLinkTitle
		dataLink := s.DataLink
		if dataLinkTitle == "" {
			dataLinkTitle = "View scrobble on Last.fm"
			dataLink = fmt.Sprintf("%vmusic/%v/%v", lfm.LastFmUrl, html.EscapeString(s.Artist), html.EscapeString(s.Name))
		}
		songButton = &client.Button{Label: dataLinkTitle, Url: dataLink}
	}

	var buttons []*client.Button
	if songButton != nil {
		buttons = []*client.Button{songButton}
	}
	if showProfile {
		buttons = []*client.Button{{Label: "Visit last.fm Profile", Url: profileUrl}}
	}

	if showProfile && songButton != nil {
		buttons = []*client.Button{{Label: "Visit last.fm Profile", Url: profileUrl}, songButton}
	}

	// Determines whether to display the heart for the smallImage
	smallUrl := "lfm_logo"
	if showLoved && s.Loved { // Change small image to heart if user enable loved and scrobble is loved
		smallUrl = "heart"
	}

	var coverUrl string
	if covers {
		coverUrl = s.CoverArtUrl
	}

	var timestamps *client.Timestamps
	if elapsed {
		timestamps = &client.Timestamps{
			Start: &s.DataTimestamp,
		}
	}

	return client.Activity{
		Details:    s.Name,
		State:      fmt.Sprintf("by %v", s.Artist),
		LargeImage: coverUrl,
		LargeText:  s.Album,
		SmallImage: smallUrl,
		SmallText:  info,
		Timestamps: timestamps,
		Buttons:    buttons,
	}

}
