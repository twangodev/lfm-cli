package main

import (
	"log/slog"
	"time"

	"github.com/lmittmann/tint"
	"github.com/xeyossr/go-discordrpc/client"
)

var loggedIn = false

// Logout closes the socket; the next Login redials to recover a dropped connection
var rpcClient = client.NewClient(discordAppId)

func getRPCLogCtx() *slog.Logger {
	return slog.With("loggedIn", loggedIn)
}

func login() {
	if loggedIn {
		return
	}
	if err := rpcClient.Login(); err != nil {
		getRPCLogCtx().Warn("Could not login to Discord.", tint.Err(err))
		logout()
		return
	}
	loggedIn = true
	getRPCLogCtx().Debug("Successfully logged into Discord's RPC Server.")
}

func logout() {
	if err := rpcClient.Logout(); err != nil {
		getRPCLogCtx().Debug("Error closing Discord RPC connection.", tint.Err(err))
	}
	loggedIn = false
	ts = time.Time{} // Force presence re-send on next login
	getRPCLogCtx().Debug("Successfully logged out of Discord's RPC Server.")
}
