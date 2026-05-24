package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/xeyossr/go-discordrpc/client"
)

var loggedIn = false

// Logout closes the socket; the next Login redials to recover a dropped connection
var rpcClient = client.NewClient(discordAppId)

func getRPCLogCtx() *log.Entry {
	return log.WithFields(log.Fields{
		"loggedIn": loggedIn,
	})
}

func login() {
	if loggedIn {
		return
	}
	if err := rpcClient.Login(); err != nil {
		getRPCLogCtx().Warnln("Could not login to Discord.")
		logout()
		return
	}
	loggedIn = true
	getRPCLogCtx().Debugln("Successfully logged into Discord's RPC Server.")
}

func logout() {
	if err := rpcClient.Logout(); err != nil {
		getRPCLogCtx().WithError(err).Debugln("Error closing Discord RPC connection.")
	}
	loggedIn = false
	getRPCLogCtx().Debugln("Successfully logged out of Discord's RPC Server.")
}
