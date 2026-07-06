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

gladia_bin() {
  if command -v "${BINARY}" >/dev/null 2>&1; then
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

maybe_prompt_completions
