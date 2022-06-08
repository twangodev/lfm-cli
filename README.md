<p align="center">
	<a href="https://last.fm/"><img alt="go-http-client" src="https://www.last.fm/static/images/lastfm_logo_facebook.15d8133be114.png" width="500"></a>
</p>

<h1 align="center">
  lfm2discord-cli
</h1>

<p align="center">
  Show your fellow gamers and friends what you're listening to on Last.FM without touching a single API Key!
</p>

<p align="center">
  <img src="https://img.shields.io/github/go-mod/go-version/lfm2discord/lfm2discord-cli">
  <img src="https://img.shields.io/github/workflow/status/lfm2discord/lfm2discord-cli/Go">
	<img src="https://img.shields.io/badge/Platforms-Windows%2C%20MacOS%2C%20Linux-orange">
  <img src="https://img.shields.io/github/license/lfm2discord/lfm2discord-cli">
</p>

<h2 align="center">
  Sample Images
</h2>
<p align="center">
  <img src="https://raw.githubusercontent.com/lfm2discord/lfm2discord-cli/master/github-assets/screenshot-1.png"><br>
  <img src="https://raw.githubusercontent.com/lfm2discord/lfm2discord-cli/master/github-assets/screenshot-2.png">
</p>

# Usage
Download the latest [release](https://github.com/lfm2discord/lfm2discord-cli/releases).

With [Discord](https://discord.com/) open, run the following binary within your console
```console
foo@bar:~$ lfm2discord-cli -u MYUSERNAME
```
For full reference on flags, run the binary with the `-h` or `--help` flag.
# Customization
Advanced users who are familiar with Go and Discord may customize the behavior of the program to their specific liking.
## Custom Discord Applications
Users may use their own discord application by specifiying their Discord application ID by modifying the `appId` constant in `main.go`. 
Rich presence assets that are used in the default application may be found in the `/discord-dev-assets/` folder.

