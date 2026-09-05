package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	lfm "github.com/twangodev/lfm-api"
	"github.com/urfave/cli/v2"
)

const name = "lfm-cli"
const version = "v1.6.3" // x-release-please-version

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

// logLevel lets main() build the handler before flags are parsed; exec() raises it for --debug.
var logLevel = new(slog.LevelVar)

func exec(ctx *cli.Context) error {

	showProfile = !ctx.Bool("hide-profile")
	showLoved = ctx.Bool("show-loved")
	covers = !ctx.Bool("rm-covers")
	elapsed = !ctx.Bool("rm-time")
	keepStatus = ctx.Bool("keep-status")
	debug = ctx.Bool("debug")
	if debug {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}

	profileUrl = fmt.Sprintf("%vuser/%v", lfm.LastFmUrl, username)

	slog.Info("Configuration loaded from arguments",
		"username", username,
		"refresh_interval", refreshInterval,
		"show_profile", showProfile,
		"show_loved", showLoved,
		"show_covers", covers,
		"show_elapsed", elapsed,
		"keep_status", keepStatus,
		"debug_enabled", debug,
	)

	for {
		slog.Debug("Cycle begin.")
		cycle()
		slog.Debug("Cycle complete.")
		time.Sleep(time.Duration(refreshInterval) * time.Second)
	}

}

func main() {

	slog.SetDefault(slog.New(tint.NewHandler(colorable.NewColorableStderr(), &tint.Options{
		Level:      logLevel,
		TimeFormat: "15:04:05.000",
		NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
	})))

	app := &cli.App{
		Name:        name,
		Description: "Show your Last.FM scrobbles on Discord!",
		Version:     version,
		Compiled:    time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "James Ding",
				Email: "james@twango.dev",
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
		slog.Error("fatal error", tint.Err(err))
		os.Exit(1)
	}

}
