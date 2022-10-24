package main

import (
	"github.com/hugolgst/rich-go/client"
	"github.com/hugolgst/rich-go/ipc"
	log "github.com/sirupsen/logrus"
)

var loggedIn = false

func getRPCLogCtx() *log.Entry {
	return log.WithFields(log.Fields{
		"loggedIn": loggedIn,
	})
}

func login() {
	getRPCLogCtx().Traceln("Attempting to close IPC Socket")
	err := ipc.CloseSocket()
	getRPCLogCtx().Warnln("IPC Socket Unable to close")
	err = client.Login(discordAppId)
	if err != nil {
		getRPCLogCtx().Warnln("Could not login to Discord.")
		logout()
	} else {
		loggedIn = true
		getRPCLogCtx().Debugln("Successfully logged into Discord's RPC Server.")
	}
}

func logout() {
	client.Logout()
	loggedIn = false
	getRPCLogCtx().Debugln("Successfully logged out of Discord's RPC Server.")
}
