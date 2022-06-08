package main

import (
	"fmt"
	"github.com/hugolgst/rich-go/client"
	"golang.org/x/net/html"
)

func createActivity(s scrobble, songLink bool) client.Activity {

	var songButton *client.Button
	// Determines whether to display profile link in buttons
	if songLink {
		dataLinkTitle := s.dataLinkTitle
		dataLink := s.dataLink
		if dataLinkTitle == "" {
			dataLinkTitle = "View scrobble on Last.fm"
			dataLink = fmt.Sprintf("%vmusic/%v/%v", lastFmUrl, html.EscapeString(s.artist), html.EscapeString(s.name))
		}
		songButton = &client.Button{Label: dataLinkTitle, Url: dataLink}
	}

	var buttons []*client.Button
	if songButton != nil {
		buttons = []*client.Button{songButton}
	}
	if showProfile {
		buttons = []*client.Button{{Label: "Visit Last.fm Profile", Url: profileUrl}}
	}

	if showProfile && songButton != nil {
		buttons = []*client.Button{{Label: "Visit Last.fm Profile", Url: profileUrl}, songButton}
	}

	// Determines whether to display the heart for the smallImage
	smallUrl := "lfm_logo"
	if showLoved && s.loved { // Change small image to heart if user enable loved and scrobble is loved
		smallUrl = "heart"
	}

	var coverUrl string
	if covers {
		coverUrl = s.coverArtUrl
	}

	var timestamps *client.Timestamps
	if elapsed {
		timestamps = &client.Timestamps{
			Start: &s.dataTimestamp,
		}
	}

	return client.Activity{
		Details:    s.name,
		State:      fmt.Sprintf("by %v", s.artist),
		LargeImage: coverUrl,
		LargeText:  s.album,
		SmallImage: smallUrl,
		SmallText:  info,
		Timestamps: timestamps,
		Buttons:    buttons,
	}

}
