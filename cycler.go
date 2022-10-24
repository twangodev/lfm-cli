package main

import (
	"fmt"
	"github.com/hugolgst/rich-go/client"
	log "github.com/sirupsen/logrus"
	"time"
)

var info = fmt.Sprintf("%v • %v", name, version)
var ts = time.Now()

func cycle() {
	s := getLatestSong(username) // Fetch latest scrobble, emptyScrobble if no new scrobble

	if keepStatus {
		login()
		if !s.active {
			err := client.SetActivity(client.Activity{
				Details:    name,
				State:      version,
				LargeImage: "lfm_logo",
			})
			if err != nil {
				log.Warnln("Failed to keep activity.")
				return
			}
		}
	} else {
		// Login logout logic
		if s.active { // Login if scrobble detected and if currently logged out
			if !loggedIn {
				log.Info("New scrobble detected. Logging in.")
				login()
			}
		} else { // No new scrobble
			if loggedIn { // Logout if logged in
				log.Info("No scrobble detected. Logging out.")
				logout()
			} else { // Retain logout state
				log.Traceln("No new scrobble detected.")
			}
			return
		}
	}

	if ts != s.dataTimestamp { // Update old timestamp to match current scrobble
		ts = s.dataTimestamp
		log.WithFields(log.Fields{"scrobbling": s}).Infoln("Updating presence.")
	} else { // Prevents update of the same scrobble, use timestamp to differentiate
		return
	}

	// First RPC attempt is without songLink
	err1 := client.SetActivity(createActivity(s, false))
	if err1 != nil {
		log.Info("Failed to set base RPC. Retrying with detailed payload.")
	} else {
		log.Traceln("Successfully set base RPC.")
	}

	// Second RPC attempt is with songLink
	err2 := client.SetActivity(createActivity(s, true))
	if err2 != nil {
		if err1 != nil {
			log.Warnln("Both attempts to set RPC failed.")
		} else {
			log.Info("Failed to set detailed RPC.")
		}
	} else {
		log.Traceln("Successfully set detailed RPC.")
	}

}
