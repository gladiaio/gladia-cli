# gladia-cli

## Direct install
Linux
```
wget https://github.com/gladiaio/gladia-cli/raw/main/dist/linux_x64_gladia && \
mv linux_x64_gladia gladia && \
chmod +x gladia
```

MacOS ARM
```
wget https://github.com/gladiaio/gladia-cli/raw/main/dist/macos_arm64_gladia && \
mv macos_arm64_gladia gladia && \
chmod +x gladia
```

## Build from source
```
$ pipenv shell
$ pip install -r requirements.txt
```

to build run
```
$ ./build.sh 
```

the resulting gladia_cli is in dist


## Usage
here is the usage:

```
$ ./gladia --help
Usage: gladia [OPTIONS]

  Transcribe an audio file using the Gladia API.

Options:
  --audio-url TEXT           URL of the audio file to be transcribed.
  --language TEXT            Language spoken in the audio file.
  --language-behaviour TEXT  Determines how to handle multi-language audio.
  --noise-reduction          Apply noise reduction to the audio.
  --output-format TEXT       Format in which to return the transcription
                             results.
  --diarization              Perform speaker diarization.
  --gladia-key TEXT          API key for Gladia. Get it at
                             https://app.gladia.io/account
  --save-gladia-key          Save the API key to a configuration file.
  --help                     Show this message and exit.
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

Transcribing audio file...
Transcript

 time_begin  time_end  probability  language  speaker        transcription
 0.09        2.07      0.49         en        not_activated  Split infinity
 2.13        5.19      0.65         en        not_activated  in a time when less is more
```

