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
gladia transcribe https://example.com/audio.mp3 -o json
```

### Output formats

| `-o` value   | Description                                      |
|--------------|--------------------------------------------------|
| `text`       | Plain transcript (default)                       |
| `json`       | Simplified JSON (utterances with timing/speaker) |
| `json-full`  | Full API response JSON                           |

### Options

```bash
gladia transcribe <source> [flags]

Flags:
  -o, --output string   Output format: text, json, json-full (default "text")
  -v, --verbose         Show progress while transcribing
      --diarize         Enable speaker diarization
      --gladia-key      API key (fallback after env and ~/.gladia)
```

## Development

```bash
make build    # build ./gladia
make dist     # cross-compile to dist/
make test     # run tests
```

## License

See [LICENSE](LICENSE).
