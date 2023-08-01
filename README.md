# gladia-cli

## Go Based CLI (New, Faster but alpha)
### Direct install
Linux X64
```
wget https://github.com/gladiaio/gladia-cli/raw/main/go/dist/gladia-darwin-amd64
```

Linux X32
```
wget https://github.com/gladiaio/gladia-cli/raw/main/go/dist/gladia-linux-armv7
```

Linux ARM
```
wget https://github.com/gladiaio/gladia-cli/raw/main/go/dist/gladia-darwin-arm64
```

MacOS Intel
```
wget https://github.com/gladiaio/gladia-cli/raw/main/go/dist/gladia-darwin-amd64
```

MacOS ARM
```
wget https://github.com/gladiaio/gladia-cli/raw/main/go/dist/gladia-darwin-arm64
```

Windows
```
wget https://github.com/gladiaio/gladia-cli/raw/main/go/dist/gladia-windows-amd64.exe
```

### Build from source
```
$ cd go
$ ./compile.sh
```

## Python Based CLI (Deprecated)

### Direct install
Linux X64
```
wget https://github.com/gladiaio/gladia-cli/raw/main/go/dist/gladia-darwin-amd64
```

Linux X32
```
wget https://github.com/gladiaio/gladia-cli/raw/main/go/dist/gladia-linux-armv7
```

Linux ARM
```
wget https://github.com/gladiaio/gladia-cli/raw/main/go/dist/gladia-darwin-arm64
```

MacOS Intel
```
wget https://github.com/gladiaio/gladia-cli/raw/main/python/dist/linux_x64_gladia && \
mv linux_x64_gladia gladia && \
chmod +x gladia
```

MacOS ARM
```
wget https://github.com/gladiaio/gladia-cli/raw/main/python/dist/macos_arm64_gladia && \
mv macos_arm64_gladia gladia && \
chmod +x gladia
```

Windows
```
wget https://github.com/gladiaio/gladia-cli/raw/main/python/dist/gladia_cli.exe
```

### Build from source
```
$ pipenv shell
$ pip install -r requirements.txt
```

to build on Macos or Linux run
```
$ ./build.sh 
```
the resulting gladia_cli is in dist 


to build on windows run
```
.\build.bat
```
the resulting gladia_cli.exe is in dist 


## Usage
here is the usage:

```
$ Usage of ./gladia:
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
```

Authentication:
1. get you Gladia key here: https://app.gladia.io/account
2. save the key if needed using
```
$ ./gladia --gladia-key MY_GLADIA_KEY --save-gladia-key
```
3. or use it inline for each request
```
$ ./gladia --gladia-key MY_GLADIA_KEY --OTHER_OPTIONS ...
```


Basic Example:
```
$ ./gladia_cli --audio-url http://files.gladia.io/example/audio-transcription/split_infinity.wav

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

