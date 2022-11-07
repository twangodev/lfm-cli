package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	lfm "github.com/twangodev/lfm-api"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

const name = "lfm-cli"
const version = "v1.1"

const discordAppId = "970003417277812736"

// Flags
var username string
var refreshInterval int
var showProfile bool
var showLoved bool
var covers bool
var elapsed bool
var keepStatus bool
var debug bool

var profileUrl string

func exec(ctx *cli.Context) error {

	showProfile = !ctx.Bool("hide-profile")
	showLoved = ctx.Bool("show-loved")
	covers = !ctx.Bool("rm-covers")
	elapsed = !ctx.Bool("rm-time")
	keepStatus = ctx.Bool("keep-status")
	debug = ctx.Bool("debug")
	if debug {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	profileUrl = fmt.Sprintf("%vuser/%v", lfm.LastFmUrl, username)

	log.WithFields(log.Fields{
		"username":         username,
		"refresh_interval": refreshInterval,
		"show_profile":     showProfile,
		"show_loved":       showLoved,
		"show_covers":      covers,
		"show_elapsed":     elapsed,
		"keep_status":      keepStatus,
		"debug_enabled":    debug,
	}).Infoln("Configuration loaded from arguments")

	for {
		log.Traceln("Cycle begin.")
		cycle()
		log.Traceln("Cycle complete.")
		time.Sleep(time.Duration(refreshInterval) * time.Second)
	}

}

func main() {

	app := &cli.App{
		Name:        name,
		Description: "Show your Last.FM scrobbles on Discord!",
		Version:     version,
		Compiled:    time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "James Ding",
				Email: "jamesding365@gmail.com",
			},
		},
		Copyright: "(c) 2022 James Ding",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "user",
				Aliases:     []string{"u"},
				Usage:       "Display Last.FM scrobbles from `USERNAME`",
				Required:    true,
				Destination: &username,
			},
			&cli.IntFlag{
				Name:        "refresh",
				Aliases:     []string{"r"},
				Usage:       "Checks Last.FM every `X` seconds for new scrobbles",
				Value:       10,
				Destination: &refreshInterval,
			},
			&cli.BoolFlag{
				Name:  "hide-profile",
				Usage: "Removes buttons to the specified Last.FM profile",
			},
			&cli.BoolFlag{
				Name:    "show-loved",
				Aliases: []string{"l"},
				Usage:   "Replaces the default smallImage key with a heart for loved songs.",
			},
			&cli.BoolFlag{
				Name:  "rm-covers",
				Usage: "Does not show album cover images.",
			},
			&cli.BoolFlag{
				Name:  "rm-time",
				Usage: "Does not show time elapsed for the scrobble.",
			},
			&cli.BoolFlag{
				Name:  "keep-status",
				Usage: "Shows status even when there is no active scrobble.",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Enable verbose and debug logging",
			},
		},
		Action: func(context *cli.Context) error {
			return exec(context)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
