# gladia-cli

## Go Based CLI (New, Faster but alpha)

### Direct install

Linux AMD (For Linux running on 64-bit AMD or Intel processors (x86_64 architecture))

```
wget https://github.com/gladiaio/gladia-cli/raw/main/dist/gladia-linux-amd64
```

Linux ARM 8 (For Linux running on 64-bit ARM processors (ARMv8 architecture)).

```
wget https://github.com/gladiaio/gladia-cli/raw/main/dist/gladia-linux-arm64
```

Linux ARM 7 (For Linux running on 32-bit ARM processors (ARMv7 architecture)).

```bash
wget https://github.com/gladiaio/gladia-cli/raw/main/dist/gladia-linux-arm7
```


MacOS Intel (For macOS running on 64-bit AMD or Intel processors (x86_64 architecture)).

```bash
wget https://github.com/gladiaio/gladia-cli/raw/main/dist/gladia-darwin-amd64
```

MacOS ARM (For macOS running on ARM64 architecture (like Apple's M1, M2 or M3 chips)).

```bash
wget https://github.com/gladiaio/gladia-cli/raw/main/dist/gladia-darwin-arm64
```

Windows (For Windows running on 64-bit AMD or Intel processors (x86_64 architecture)).

```bash
wget https://github.com/gladiaio/gladia-cli/raw/main/dist/gladia-windows-amd64.exe
```

### Build from source

```bash
make build
```

## Usage

here is the usage:

```bash
Usage of ./gladia:
  -audio-file string
        Path to the audio file
  -audio-url string
        URL of the audio file
  -diarization
        Enable diarization
  -diarization-max-speakers int
        Maximum number of speakers for diarization
  -direct-translate
        Enable direct translation
  -direct-translate-language string
        Language for direct translation
  -gladia-key string
        Gladia API key
  -language string
        Language for transcription (default "english")
  -language-behaviour string
        Language behavior (manual, automatic single language, automatic multiple languages) (default "automatic multiple languages")
  -noise-reduction
        Enable noise reduction
  -output-format string
        Output format (table, csv, json, srt, vtt, txt) (default "table")
  -save-gladia-key
        Save Gladia API key
  -transcription-hint string
        Transcription hint
  -transcription-language-list
        List available languages for transcription
  -translation-language-list
        List available languages for translation
  -verbose
        Enable verbose printing (default=true)
```

Authentication:

1.  get you Gladia key here: https://app.gladia.io/account
2.  save the key if needed using
3.  or use it inline for each request

Basic Example:

```bash
./gladia_cli --audio-url http://files.gladia.io/example/audio-transcription/split_infinity.wav

+------------+----------+----------+-----------------------+--------------------------------+
| TIME BEGIN | TIME END | LANGUAGE |        SPEAKER        |         TRANSCRIPTION          |
+------------+----------+----------+-----------------------+--------------------------------+
|       0.18 |     4.68 | en       | speaker_not_activated | Split infinity in a time when  |
|            |          |          |                       | less is more,                  |
|       5.52 |     7.76 | en       | speaker_not_activated | where too much is never        |
|            |          |          |                       | enough.                        |
|       8.51 |    10.79 | en       | speaker_not_activated | There is always hope for the   |
|            |          |          |                       | future.                        |
|      11.71 |    14.11 | en       | speaker_not_activated | The future can be read from    |
|            |          |          |                       | the past.                      |
|      14.57 |    19.91 | en       | speaker_not_activated | The past foreshadows the       |
|            |          |          |                       | present and the present hasn't |
|            |          |          |                       | been written yet.              |
+------------+----------+----------+-----------------------+--------------------------------+
```

```bash
./gladia --gladia-key MY_GLADIA_KEY --OTHER_OPTIONS ...
```

```bash
./gladia --gladia-key MY_GLADIA_KEY --save-gladia-key
```
