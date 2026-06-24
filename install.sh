#!/bin/sh
set -eu

REPO="gladiaio/gladia-cli"
BINARY="gladia"

detect_os() {
  case "$(uname -s)" in
    Darwin) printf '%s\n' darwin ;;
    Linux) printf '%s\n' linux ;;
    *)
      echo "error: unsupported OS: $(uname -s) (macOS and Linux only)" >&2
      exit 1
      ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64 | amd64) printf '%s\n' amd64 ;;
    aarch64 | arm64) printf '%s\n' arm64 ;;
    armv7l | armv6l) printf '%s\n' armv7 ;;
    i386 | i686) printf '%s\n' 386 ;;
    *)
      echo "error: unsupported architecture: $(uname -m)" >&2
      exit 1
      ;;
  esac
}

fetch_latest_tag() {
  if [ -n "${GITHUB_TOKEN:-}" ]; then
    curl -fsSL \
      -H "Authorization: Bearer ${GITHUB_TOKEN}" \
      -H "Accept: application/vnd.github+json" \
      "https://api.github.com/repos/${REPO}/releases/latest" \
      | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
      | head -n 1
  else
    curl -fsSL \
      -H "Accept: application/vnd.github+json" \
      "https://api.github.com/repos/${REPO}/releases/latest" \
      | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
      | head -n 1
  fi
}

install_binary() {
  src="$1"
  dest_dir="$2"

  if [ -w "$dest_dir" ]; then
    install -m 755 "$src" "$dest_dir/"
  else
    if ! command -v sudo >/dev/null 2>&1; then
      echo "error: cannot write to ${dest_dir} and sudo is not available" >&2
      echo "hint: set GLADIA_INSTALL_DIR to a writable directory (e.g. \$HOME/.local/bin)" >&2
      exit 1
    fi
    sudo install -m 755 "$src" "$dest_dir/"
  fi
}

os="$(detect_os)"
arch="$(detect_arch)"
tag="$(fetch_latest_tag)"

if [ -z "$tag" ]; then
  echo "error: could not determine latest release" >&2
  exit 1
fi

version="${tag#v}"
archive="${BINARY}_${version}_${os}_${arch}.tar.gz"
url="https://github.com/${REPO}/releases/download/${tag}/${archive}"
install_dir="${GLADIA_INSTALL_DIR:-/usr/local/bin}"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT INT HUP TERM

echo "Installing ${BINARY} ${tag} (${os}/${arch})..."

if ! curl -fsSL "$url" | tar -xz -C "$tmpdir"; then
  echo "error: failed to download ${url}" >&2
  exit 1
fi

if [ ! -f "${tmpdir}/${BINARY}" ]; then
  echo "error: ${BINARY} binary not found in archive" >&2
  exit 1
fi

install_binary "${tmpdir}/${BINARY}" "$install_dir"

if command -v "${BINARY}" >/dev/null 2>&1; then
  echo "Installed ${BINARY} $(command -v "${BINARY}") ($("${BINARY}" --version 2>/dev/null || true))"
else
  echo "Installed ${BINARY} to ${install_dir}/${BINARY}"
  echo "Ensure ${install_dir} is in your PATH"
fi
