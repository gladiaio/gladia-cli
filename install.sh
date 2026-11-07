#!/bin/sh
# Gladia CLI installer for macOS and Linux.
# Downloads the latest GitHub release, verifies its SHA-256 against the
# release's checksums.txt, and installs the binary.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/gladiaio/gladia-cli/main/install.sh | sh
#
# Environment variables:
#   GLADIA_INSTALL_DIR           Install directory (default: /usr/local/bin)
#   GITHUB_TOKEN                 Optional, for GitHub API rate limits (only sent to api.github.com)
#   GLADIA_NO_COMPLETION_PROMPT  Skip the shell completion prompt
#   GLADIA_INSTALL_API_URL       Override the "latest release" API URL (testing/mirrors)
#   GLADIA_INSTALL_DOWNLOAD_BASE Override the release download base URL (testing/mirrors)
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
  default_api_url="https://api.github.com/repos/${REPO}/releases/latest"
  api_url="${GLADIA_INSTALL_API_URL:-$default_api_url}"
  # Only send GITHUB_TOKEN to the real GitHub API, never to an override URL.
  if [ -n "${GITHUB_TOKEN:-}" ] && [ "$api_url" = "$default_api_url" ]; then
    curl -fsSL \
      -H "Authorization: Bearer ${GITHUB_TOKEN}" \
      -H "Accept: application/vnd.github+json" \
      "${api_url}" \
      | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
      | head -n 1
  else
    curl -fsSL \
      -H "Accept: application/vnd.github+json" \
      "${api_url}" \
      | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
      | head -n 1
  fi
}

# Verify a downloaded archive against checksums.txt (published per release by
# GoReleaser). Refuses to continue if the entry is missing or the hash differs.
verify_checksum() {
  archive_path="$1"
  checksums_file="$2"
  archive_base="$(basename "$archive_path")"

  expected="$(awk -v f="$archive_base" '$2 == f { print $1; exit }' "$checksums_file")"
  if [ -z "$expected" ]; then
    echo "error: no sha256 entry for ${archive_base} in checksums.txt" >&2
    return 1
  fi

  if command -v sha256sum >/dev/null 2>&1; then
    actual="$(sha256sum "$archive_path" | awk '{ print $1 }')"
  elif command -v shasum >/dev/null 2>&1; then
    actual="$(shasum -a 256 "$archive_path" | awk '{ print $1 }')"
  else
    echo "error: sha256sum/shasum not found; cannot verify checksum" >&2
    return 1
  fi

  if [ "$actual" != "$expected" ]; then
    echo "error: checksum mismatch for ${archive_base}" >&2
    echo "  expected: ${expected}" >&2
    echo "  actual:   ${actual}" >&2
    return 1
  fi
  echo "Checksum verified: ${archive_base}"
}

# Defense-in-depth against zip-slip style archives: GNU tar and bsdtar already
# strip leading '/' and refuse '..' members (exiting non-zero), but behavior
# varies across tar flavors, so reject unsafe archives explicitly. Returns 0 if
# the archive contains any entry with an absolute path, a '..' component, or a
# symlink/hardlink member.
has_unsafe_tar_entry() {
  if tar -tzf "$1" | grep -Eq '(^|/)\.\.(/|$)|^/'; then
    return 0
  fi
  if tar -tvzf "$1" | grep -Eq '^[lh]'; then
    return 0
  fi
  return 1
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

gladia_bin() {
  # Prefer the binary we just installed over any other 'gladia' on PATH, so a
  # stale or malicious earlier PATH entry cannot shadow the release we fetched.
  if [ -x "${install_dir}/${BINARY}" ]; then
    printf '%s\n' "${install_dir}/${BINARY}"
  elif command -v "${BINARY}" >/dev/null 2>&1; then
    command -v "${BINARY}"
  else
    printf '%s\n' "${install_dir}/${BINARY}"
  fi
}

print_completion_hints() {
  echo ""
  echo "Shell tab completion is available. To set up manually:"
  echo "  gladia completion --help"
  echo "  gladia completion bash   # requires bash-completion package"
  echo "  gladia completion zsh"
  echo "  gladia completion fish"
}

install_fish_completions() {
  gladia_cmd="$(gladia_bin)"
  mkdir -p "${HOME}/.config/fish/completions"
  "${gladia_cmd}" completion fish > "${HOME}/.config/fish/completions/gladia.fish"
  echo "Wrote fish completion to ~/.config/fish/completions/gladia.fish"
  echo "Restart your shell for completion to take effect."
}

install_zsh_completions() {
  gladia_cmd="$(gladia_bin)"
  mkdir -p "${HOME}/.zsh/completions"
  "${gladia_cmd}" completion zsh > "${HOME}/.zsh/completions/_gladia"

  zshrc="${ZDOTDIR:-${HOME}}/.zshrc"
  if [ -f "${zshrc}" ] && grep -q '# gladia-cli completions' "${zshrc}" 2>/dev/null; then
    echo "Updated zsh completion at ~/.zsh/completions/_gladia"
    echo "Restart your shell for completion to take effect."
    return 0
  fi

  if [ ! -f "${zshrc}" ]; then
    : > "${zshrc}"
  fi

  cat >> "${zshrc}" <<'EOF'

# gladia-cli completions
fpath=(~/.zsh/completions $fpath)
autoload -U compinit; compinit
EOF
  echo "Wrote zsh completion to ~/.zsh/completions/_gladia"
  echo "Restart your shell for completion to take effect."
}

install_bash_completions() {
  gladia_cmd="$(gladia_bin)"
  target=""

  if [ -d "${HOME}/.local/share/bash-completion" ]; then
    mkdir -p "${HOME}/.local/share/bash-completion/completions"
    target="${HOME}/.local/share/bash-completion/completions/gladia"
  elif command -v brew >/dev/null 2>&1; then
    brew_prefix="$(brew --prefix 2>/dev/null || true)"
    if [ -n "${brew_prefix}" ] && [ -d "${brew_prefix}/etc/bash_completion.d" ]; then
      target="${brew_prefix}/etc/bash_completion.d/gladia"
    fi
  fi

  if [ -n "${target}" ]; then
    "${gladia_cmd}" completion bash > "${target}"
    echo "Wrote bash completion to ${target}"
    echo "Restart your shell for completion to take effect."
    return 0
  fi

  echo "Could not find a bash-completion directory."
  echo "Install the bash-completion package, then run:"
  echo "  gladia completion bash > ~/.local/share/bash-completion/completions/gladia"
  echo "Or add to your shell rc:"
  echo "  source <(gladia completion bash)"
}

install_completions() {
  shell_name="$(basename "${SHELL:-}")"
  case "${shell_name}" in
    fish)
      install_fish_completions
      ;;
    zsh)
      install_zsh_completions
      ;;
    bash)
      install_bash_completions
      ;;
    *)
      echo "Unknown shell (${shell_name}). Run: gladia completion --help"
      print_completion_hints
      ;;
  esac
}

maybe_prompt_completions() {
  if [ -n "${GLADIA_NO_COMPLETION_PROMPT:-}" ]; then
    print_completion_hints
    return 0
  fi
  if [ ! -t 0 ] || [ ! -t 1 ]; then
    print_completion_hints
    return 0
  fi

  printf "Install shell tab completion? [y/N] "
  reply=""
  read -r reply || true
  case "${reply}" in
    [yY]|[yY][eE][sS])
      install_completions
      ;;
    *)
      print_completion_hints
      ;;
  esac
}

os="$(detect_os)"
arch="$(detect_arch)"
tag="$(fetch_latest_tag)"

if [ -z "$tag" ]; then
  echo "error: could not determine latest release" >&2
  exit 1
fi

# Sanity-check the tag before interpolating it into URLs or filenames.
case "$tag" in
  [!A-Za-z0-9]* | *[!A-Za-z0-9._+-]*)
    echo "error: unexpected tag format from GitHub API: ${tag}" >&2
    exit 1
    ;;
esac

version="${tag#v}"
archive="${BINARY}_${version}_${os}_${arch}.tar.gz"
download_base="${GLADIA_INSTALL_DOWNLOAD_BASE:-https://github.com/${REPO}/releases/download}"
download_base="${download_base%/}"
url="${download_base}/${tag}/${archive}"
checksums_url="${download_base}/${tag}/checksums.txt"
install_dir="${GLADIA_INSTALL_DIR:-/usr/local/bin}"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT INT HUP TERM

echo "Installing ${BINARY} ${tag} (${os}/${arch})..."

archive_path="${tmpdir}/${archive}"
if ! curl -fsSL "$url" -o "$archive_path"; then
  echo "error: failed to download ${url}" >&2
  exit 1
fi
if ! curl -fsSL "$checksums_url" -o "${tmpdir}/checksums.txt"; then
  echo "error: failed to download ${checksums_url}" >&2
  exit 1
fi
if ! verify_checksum "$archive_path" "${tmpdir}/checksums.txt"; then
  exit 1
fi
if has_unsafe_tar_entry "$archive_path"; then
  echo "error: archive contains unsafe entries (absolute path, '..' or link); aborting" >&2
  exit 1
fi
if ! tar -xzf "$archive_path" -C "$tmpdir"; then
  echo "error: failed to extract ${archive}" >&2
  exit 1
fi

if [ ! -f "${tmpdir}/${BINARY}" ]; then
  echo "error: ${BINARY} binary not found in archive" >&2
  exit 1
fi

install_binary "${tmpdir}/${BINARY}" "$install_dir"

# Run the freshly installed binary directly; resolving via PATH here would
# let a stale/malicious 'gladia' earlier in PATH shadow the new release.
installed_bin="${install_dir}/${BINARY}"
if [ -x "$installed_bin" ]; then
  echo "Installed ${BINARY} to ${installed_bin} ($("${installed_bin}" --version 2>/dev/null || true))"
else
  echo "Installed ${BINARY} to ${installed_bin}"
  echo "Ensure ${install_dir} is in your PATH"
fi

maybe_prompt_completions
