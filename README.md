# gladia-cli

Transcribe audio files with the [Gladia](https://www.gladia.io/) pre-recorded API (v2).

## Install

```bash
make build   # → ./gladia
```

Or download a binary from [GitHub releases](https://github.com/gladiaio/gladia-cli/releases).

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
```

## Commands

| Command | Description |
|---------|-------------|
| `transcribe <source>` | Transcribe an audio |
| `auth set <key>` | Save API key to `~/.gladia` |
| `languages` | List supported ISO 639-1 codes |

## Flags (`transcribe`)

| Flag | Default | Description |
|------|---------|-------------|
| `-o`, `--output` | `text` | Output: `text`, `json`, `json-full`, `srt`, `vtt` |
| `--language` | — | Expected language(s), comma-separated (`en` or `en,fr,de`) |
| `--code-switching`, `--code-switch` | off | Detect language per utterance |
| `--diarize` | off | **Optional.** Identify speakers in the transcript |
| `-v`, `--verbose` | off | Show progress while polling |

**Global flag** (any command): `--gladia-key` — API key if not in env or `~/.gladia`

## Language

| Goal | What to run |
|------|-------------|
| Auto-detect | `transcribe <source>` |
| Constrain detection | `--language en,fr,de` (no code switching) |
| Code switching | `--code-switching` (+ optional `--language` hints) |

- **`--language`** — tells Gladia which language(s) to expect. Several codes (`en,fr,de`) narrow detection; they do **not** turn on code switching.
- **`--code-switching`** — separate option: re-detect language on each utterance. Combine with `--language` when you know which languages may appear.

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
