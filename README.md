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
| Auto-detect | Omit `--language` |
| Single language | `--language en` |
| **Code switching** (2+ languages in the same audio) | `--language en,fr,de` |

For **code switching**, pass **several comma-separated ISO codes** (2–5 recommended, e.g. `en,fr,de`). The CLI sends `language_config.languages` and turns on `code_switching` automatically — you do **not** need `--code-switching` unless you want to be explicit:

```bash
gladia transcribe interview.mp3 --language en,fr,de
# equivalent:
gladia transcribe interview.mp3 --language en,fr,de --code-switching
```

`--code-switch` is an alias for `--code-switching`. Both require **at least two** languages in `--language`.

### Options

```bash
gladia transcribe <source> [flags]

Flags:
  -o, --output string       Output format: text, json, json-full, srt, vtt (default "text")
      --language string     ISO codes: en (single) or en,fr,de (code switching)
      --code-switching      Force code switching (needs 2+ languages in --language)
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
