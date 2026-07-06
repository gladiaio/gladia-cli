# gladia-cli

## Install

### macOS & Linux

```bash
curl -fsSL https://github.com/gladiaio/gladia-cli/releases/latest/download/install.sh | sh
```

### Windows

```powershell
powershell -c "irm https://github.com/gladiaio/gladia-cli/releases/latest/download/install.ps1 | iex"
```

Other platforms: [GitHub releases](https://github.com/gladiaio/gladia-cli/releases).

The installer prompts to set up shell tab completion when run interactively. To skip the prompt (e.g. in CI), set `GLADIA_NO_COMPLETION_PROMPT=1`.

## Shell completion

When you install via `install.sh` or `install.ps1`, the script asks whether to configure tab completion for your shell. You can also set it up manually:

```bash
gladia completion --help
```

### bash

Requires the [bash-completion](https://github.com/scop/bash-completion) package (on macOS: `brew install bash-completion@2`).

```bash
# current session
source <(gladia completion bash)

# persistent (user directory)
mkdir -p ~/.local/share/bash-completion/completions
gladia completion bash > ~/.local/share/bash-completion/completions/gladia
```

### zsh

```bash
mkdir -p ~/.zsh/completions
gladia completion zsh > ~/.zsh/completions/_gladia

# add to ~/.zshrc if not already present:
# fpath=(~/.zsh/completions $fpath)
# autoload -U compinit; compinit
```

### fish

```bash
mkdir -p ~/.config/fish/completions
gladia completion fish > ~/.config/fish/completions/gladia.fish
```

### PowerShell

```powershell
gladia completion powershell | Out-File -Append -Encoding utf8 $PROFILE
```

Restart your shell after installing completions.

## Install (from source)

```bash
make build   # → ./gladia
```

## Setup

Get an API key at [app.gladia.io/account](https://app.gladia.io/account), then either:

```bash
export GLADIA_API_KEY=your_key
# or
./gladia auth set your_key   # saves to ~/.gladia (mode 0600)
```

**Credential order:** `GLADIA_API_KEY` → `~/.gladia` → `--gladia-key`

## Usage

```bash
./gladia transcribe <file-or-url> [flags]
```

**Examples**

```bash
./gladia transcribe meeting.wav
./gladia transcribe https://example.com/audio.mp3 -o json
./gladia transcribe podcast.mp3 --language en,fr,de
./gladia transcribe mixed.mp3 --code-switching --language en,fr
./gladia transcribe call.wav --diarize -o srt
./gladia transcribe podcast.mp3 --model solaria-3 --language en
```

## Commands

| Command               | Description                                                 |
| --------------------- | ----------------------------------------------------------- |
| `transcribe <source>` | Transcribe an audio                                         |
| `auth set <key>`      | Save API key to `~/.gladia`                                 |
| `languages`           | List supported ISO 639-1 codes                              |
| `completion <shell>`  | Generate shell tab completion (bash, zsh, fish, powershell) |

## Flags (`transcribe`)

| Flag                       | Default | Description                                                                                                                                              |
| -------------------------- | ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `-o`, `--output`           | `text`  | Output: `text`, `json`, `json-full`, `srt`, `vtt`                                                                                                        |
| `--language`               | —       | Expected language(s), comma-separated (`en` or `en,fr,de`); narrows detection, does not enable code switching                                            |
| `--cs`, `--code-switching` | off     | Re-detect language on each utterance (mixed-language audio; solaria-1 only)                                                                              |
| `--diarize`                | off     | **Optional.** Identify speakers in the transcript                                                                                                        |
| `--model`                  | —       | STT model: `solaria-1` or `solaria-3`. Solaria-3 accepts at most one `--language` (`en`, `fr`, `de`, `es`, or `it`) and does not support code switching. |
| `-v`, `--verbose`          | off     | Show progress while polling                                                                                                                              |

**Global flag** (any command): `--gladia-key` — API key if not in env or `~/.gladia`

## Language

| Goal                | What to run                                                  |
| ------------------- | ------------------------------------------------------------ |
| Auto-detect         | `transcribe <source>`                                        |
| Constrain detection | `--language en,fr,de` (no code switching)                    |
| Code switching      | `--cs` or `--code-switching` (+ optional `--language` hints) |

- **`--language`** — limits which language(s) Gladia considers (`en,fr,de` is a hint list, not per-utterance switching).
- **`--cs`** / **`--code-switching`** — turns on per-utterance language detection. Add `--language` to restrict which languages may appear. Not available with `solaria-3`.

```bash
./gladia languages   # list valid codes
```

## Diarization (optional)

Use **`--diarize`** when you need **who spoke when**. Off by default.

- Works with any output format; most useful with `-o text`, `srt`, or `vtt`.
- Speaker labels are included in the output (e.g. `Speaker 0: …`).

```bash
./gladia transcribe meeting.wav --diarize
./gladia transcribe panel.mp3 --diarize -o srt
```

## Develop

```bash
make build && make test && make dist
```

## License

[MIT](LICENSE)
