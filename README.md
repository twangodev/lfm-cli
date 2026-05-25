<p align="center">
	<a href="https://last.fm/"><img alt="go-http-client" src="https://www.last.fm/static/images/lastfm_logo_facebook.15d8133be114.png" width="500"></a>
</p>

<h1 align="center">
  lfm-cli
</h1>

<p align="center">
  Show your fellow gamers and friends what you're listening to on Last.FM without touching a single API Key!
</p>

<p align="center">
  <img src="https://img.shields.io/github/go-mod/go-version/twangodev/lfm-cli">
  <img alt="GitHub all releases" src="https://img.shields.io/github/downloads/twangodev/lfm-cli/total">
  <img src="https://img.shields.io/github/actions/workflow/status/twangodev/lfm-cli/go.yml?branch=main">
  <img src="https://img.shields.io/badge/Platforms-Windows%2C%20MacOS%2C%20Linux-orange">
  <img src="https://img.shields.io/github/license/twangodev/lfm-cli">
</p>

<h2 align="center">
  Sample Images
</h2>
<p align="center">
  <img src=".github/assets/screenshot-1.png"><br>
  <img src=".github/assets/screenshot-2.png">
</p>

# Getting Started

lfm-cli works right out of the box - no configuration needed.

## Installation

**Homebrew (macOS):**
```console
brew install twangodev/tap/lfm-cli
```

**Scoop (Windows):**
```console
scoop bucket add twangodev https://github.com/twangodev/scoop-bucket
scoop install lfm-cli
```

**Install script (macOS / Linux / FreeBSD):**
```console
curl -fsSL https://twango.dev/install/lfm-cli | sh
```

**Linux packages:** download the `.deb`, `.rpm`, or `.apk` from the
[latest release](https://github.com/twangodev/lfm-cli/releases/latest).

**Manual:** download the archive for your platform from the
[releases page](https://github.com/twangodev/lfm-cli/releases). Releases include
SBOMs and cosign-signed checksums generated on GitHub Actions.

## Usage

**With [Discord](https://discord.com/) open**, run the following binary in your console

```console
foo@bar:~$ lfm-cli -u MYUSERNAME
```

For full reference on flags, run the binary with the `-h` or `--help` flag.
