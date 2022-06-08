package main

import (
	"github.com/hugolgst/rich-go/client"
	log "github.com/sirupsen/logrus"
)

var loggedIn = false

func getRPCLogCtx() *log.Entry {
	return log.WithFields(log.Fields{
		"loggedIn": loggedIn,
	})
}

func login() {
	err := client.Login(discordAppId)
	if err != nil {
		getRPCLogCtx().Warnln("Could not login to Discord.")
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
