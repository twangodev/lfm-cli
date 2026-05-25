#!/bin/sh
# lfm-cli installer. Usage:
#   curl -fsSL https://twango.dev/install/lfm-cli | sh
# Env overrides: VERSION (e.g. v1.5.0), INSTALL_DIR.
set -eu

REPO="twangodev/lfm-cli"
BIN="lfm-cli"

err() { printf 'error: %s\n' "$1" >&2; exit 1; }

# --- detect OS ---
os=$(uname -s)
case "$os" in
	Darwin) os="darwin" ;;
	Linux) os="linux" ;;
	FreeBSD) os="freebsd" ;;
	MINGW* | MSYS* | CYGWIN* | Windows_NT)
		err "Windows is not supported by this script. Use Scoop: scoop install lfm-cli" ;;
	*) err "unsupported OS: $os" ;;
esac

# --- detect arch (must match GoReleaser archive names) ---
arch=$(uname -m)
case "$arch" in
	x86_64 | amd64) arch="amd64" ;;
	arm64 | aarch64) arch="arm64" ;;
	armv7l | armv6l | arm) arch="arm" ;;
	i386 | i686) arch="386" ;;
	*) err "unsupported architecture: $arch" ;;
esac

# --- resolve version ---
version="${VERSION:-}"
if [ -z "$version" ]; then
	version=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" |
		grep '"tag_name"' | head -n1 | cut -d'"' -f4)
fi
[ -n "$version" ] || err "could not resolve latest release version"

# --- download + verify ---
archive="${BIN}_${version#v}_${os}_${arch}.tar.gz"
base="https://github.com/${REPO}/releases/download/${version}"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

printf 'Downloading %s %s (%s/%s)...\n' "$BIN" "$version" "$os" "$arch"
curl -fsSL "${base}/${archive}" -o "${tmp}/${archive}" ||
	err "download failed: ${base}/${archive}"
curl -fsSL "${base}/checksums.txt" -o "${tmp}/checksums.txt" ||
	err "checksums download failed"

# verify SHA256 (sha256sum on Linux, shasum on macOS/BSD)
expected=$(grep " ${archive}\$" "${tmp}/checksums.txt" | cut -d' ' -f1)
[ -n "$expected" ] || err "no checksum entry for ${archive}"
if command -v sha256sum >/dev/null 2>&1; then
	actual=$(sha256sum "${tmp}/${archive}" | cut -d' ' -f1)
else
	actual=$(shasum -a 256 "${tmp}/${archive}" | cut -d' ' -f1)
fi
[ "$expected" = "$actual" ] || err "checksum mismatch for ${archive}"

# --- extract + install ---
tar -xzf "${tmp}/${archive}" -C "$tmp"
[ -f "${tmp}/${BIN}" ] || err "binary ${BIN} not found in archive"

dir="${INSTALL_DIR:-}"
if [ -z "$dir" ]; then
	if [ -w /usr/local/bin ] 2>/dev/null; then dir="/usr/local/bin"; else dir="${HOME}/.local/bin"; fi
fi
mkdir -p "$dir"
install -m 0755 "${tmp}/${BIN}" "${dir}/${BIN}"

printf 'Installed %s to %s/%s\n' "$BIN" "$dir" "$BIN"
case ":${PATH}:" in
	*":${dir}:"*) ;;
	*) printf "Note: %s is not on your PATH. Add it:\n  export PATH=\"%s:\$PATH\"\n" "$dir" "$dir" ;;
esac
"${dir}/${BIN}" --version || true
