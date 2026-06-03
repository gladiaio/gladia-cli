# gladia-cli

Command-line tool for Gladia pre-recorded speech-to-text (v2 API).

## Install

### Pre-built binaries

Download from [GitHub releases](https://github.com/gladiaio/gladia-cli/releases) or build from source:

```bash
make build
```

The binary is written to `./gladia`.

## Authentication

Credentials are resolved in this order:

1. `GLADIA_API_KEY` environment variable
2. `~/.gladia` file
3. `--gladia-key` flag

Save a key to disk:

```bash
gladia auth set YOUR_API_KEY
```

Get a key at [app.gladia.io/account](https://app.gladia.io/account).

## Quickstart

```bash
export GLADIA_API_KEY=your_key

gladia transcribe meeting.wav
gladia transcribe audio.mp3 -o text
gladia transcribe podcast.mp3 --language en
gladia transcribe interview.mp3 --language en,fr,de
gladia transcribe call.wav --diarize -o srt
gladia transcribe https://example.com/audio.mp3 -o json
```

List supported language codes:

```bash
gladia languages
```

### Output formats

| `-o` value   | Description                                      |
|--------------|--------------------------------------------------|
| `text`       | Plain transcript (default)                       |
| `json`       | Simplified JSON (utterances with timing/speaker) |
| `json-full`  | Full API response JSON                           |
| `srt`        | SubRip subtitles (from utterances)                 |
| `vtt`        | WebVTT subtitles (from utterances)               |

With `--diarize`, `text`, `srt`, and `vtt` include speaker labels.

### Language & code switching

| Scenario | Command |
|----------|---------|
| Auto-detect | Omit `--language` and `--code-switching` |
| Single language | `--language en` |
| Code switching (no language hint) | `--code-switching` |
| Code switching + one language | `--code-switching --language en` |
| Code switching + several languages | `--language en,fr,de` (or add `--code-switching`) |

`--language` is optional with code switching. You can pass zero, one, or several comma-separated ISO codes as hints. Listing **2–5 expected codes** (e.g. `en,fr,de`) improves accuracy; with multiple codes, code switching is turned on automatically.

```bash
gladia transcribe interview.mp3 --code-switching
gladia transcribe interview.mp3 --code-switching --language en
gladia transcribe interview.mp3 --language en,fr,de
```

`--code-switch` is an alias for `--code-switching`.

### Options

```bash
gladia transcribe <source> [flags]

Flags:
  -o, --output string       Output format: text, json, json-full, srt, vtt (default "text")
      --language string     Optional ISO codes, comma-separated (e.g. en or en,fr,de)
      --code-switching      Detect language per utterance
      --code-switch         Alias for --code-switching
  -v, --verbose             Show progress while transcribing
      --diarize             Enable speaker diarization
      --gladia-key          API key (fallback after env and ~/.gladia)
```

## Development

```bash
make build    # build ./gladia
make dist     # cross-compile to dist/
make test     # run tests
```

## License

See [LICENSE](LICENSE).
