package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/lmittmann/tint"
	lfm "github.com/twangodev/lfm-api"
	"github.com/xeyossr/go-discordrpc/client"
)

var info = fmt.Sprintf("%v • %v", name, version)
var ts = time.Now()

// Swappable in tests
var getActiveScrobble = lfm.GetActiveScrobble

func cycle() {
	s, err := getActiveScrobble(username) // Fetch latest scrobble, emptyScrobble if no new scrobble
	if err != nil {                       // Transient fetch failure is not "stopped scrobbling"
		slog.Warn("Could not fetch scrobble state. Retaining presence until next cycle.", tint.Err(err))
		return
	}

	if keepStatus {
		login()
		if !s.Active {
			err := rpcClient.SetActivity(client.Activity{
				Details:    name,
				State:      version,
				LargeImage: "lfm_logo",
			})
			if err != nil {
				slog.Warn("Failed to keep activity. Dropping connection to reconnect next cycle.", tint.Err(err))
				logout()
				return
			}
		}
	} else {
		// Login logout logic
		if s.Active { // Login if scrobble detected and if currently logged out
			if !loggedIn {
				slog.Info("New scrobble detected. Logging in.")
				login()
			}
		} else { // No new scrobble
			if loggedIn { // Logout if logged in
				slog.Info("No scrobble detected. Logging out.")
				logout()
			} else { // Retain logout state
				slog.Debug("No new scrobble detected.")
			}
			return
		}
	}

	if ts != s.DataTimestamp { // Update old timestamp to match current scrobble
		ts = s.DataTimestamp
		slog.Info("Updating presence.", "scrobbling", s)
	} else { // Prevents update of the same scrobble, use timestamp to differentiate
		return
	}

	// First RPC attempt is without songLink
	err1 := rpcClient.SetActivity(createActivity(s, false))
	if err1 != nil {
		slog.Info("Failed to set base RPC. Retrying with detailed payload.", tint.Err(err1))
	} else {
		slog.Debug("Successfully set base RPC.")
	}

	// Second RPC attempt is with songLink
	err2 := rpcClient.SetActivity(createActivity(s, true))
	if err2 != nil {
		if err1 != nil {
			slog.Warn("Both attempts to set RPC failed. Reconnecting next cycle.", "base_err", err1, "detailed_err", err2)
			logout()
		} else {
			slog.Info("Failed to set detailed RPC.", tint.Err(err2))
		}
	} else {
		slog.Debug("Successfully set detailed RPC.")
	}
}
