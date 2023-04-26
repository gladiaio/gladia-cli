import click
import requests
import os
import magic
from prettytable import PrettyTable

GLADIA_API_URL = "https://api.gladia.io/audio/text/audio-transcription/"
CONFIG_PATH = os.path.join(os.path.expanduser("~"), ".gladia")

class Color:
   PURPLE = '\033[95m'
   CYAN = '\033[96m'
   DARKCYAN = '\033[36m'
   BLUE = '\033[94m'
   GREEN = '\033[92m'
   YELLOW = '\033[93m'
   RED = '\033[91m'
   BOLD = '\033[1m'
   UNDERLINE = '\033[4m'
   END = '\033[0m'

def save_gladia_key_to_file(gladia_key):
    with open(CONFIG_PATH, "w") as f:
        f.write(gladia_key)
    click.echo("Gladia API key saved to {}".format(CONFIG_PATH))

def get_gladia_key():
    try:
        with open(CONFIG_PATH, "r") as f:
            return f.read().strip()
    except FileNotFoundError:
        click.echo("please provide your gladia key using --gladia-key or --save-gladia-key")
        return None


@click.command()
@click.option("--audio-url", help="URL of the audio file to be transcribed.")
@click.option("--audio-file", help="Path to the audio file to be transcribed.")
@click.option("--language-behaviour", default="automatic multiple languages", help="Determines how to handle multi-language audio.")
@click.option("--language", default="english", help="Language spoken in the audio file.")
@click.option("--transcription-hint", default="general", help="Hint to the transcription model. You can pass names, topics, custom vocabulary, etc.")
@click.option("--noise-reduction", is_flag=True, help="Apply noise reduction to the audio.")
@click.option("--diarization", is_flag=True, help="Perform speaker diarization.")
@click.option("--diarization-max-speakers", default="3", help="Determines the maximum number of speakers to be detected.")
@click.option("--direct-translate", is_flag=True, help="Activate direct translation to the specified language.")
@click.option("--direct-translate-language", default="english", help="Language to which to translate the transcription, need to activate the direct translation using --direct-translate.")
@click.option("--text-emotion", is_flag=True, help="Activate text emotion recognition.")
@click.option("--summarization", is_flag=True, help="Activate summarization.")
@click.option("--output-format", default="table", help="Format in which to return the transcription results. Possible values: table, json, text, srt, vtt, plain.")
@click.option("--gladia-key", help="API key for Gladia. Get it at https://app.gladia.io/account")
@click.option("--save-gladia-key", is_flag=True, help="Save the API key to a configuration file.")
def transcribe(
    audio_url: str,
    audio_file: str,
    language_behaviour: str,
    language: str,
    transcription_hint: str,
    noise_reduction: bool,
    diarization: bool,
    diarization_max_speakers: int,
    direct_translate: bool,
    direct_translate_language: str,
    text_emotion: bool,
    summarization: bool,
    output_format: str,
    gladia_key: str,
    save_gladia_key: bool
    ):
    """
    Transcribe an audio file or an audio url using the Gladia API.
    """
    if gladia_key is None:
        gladia_key = get_gladia_key()

    if save_gladia_key is True:
        save_gladia_key_to_file(gladia_key)
    
    if gladia_key is None and not save_gladia_key:
        click.echo("Error: Gladia API key not found.")
        return
    
    if save_gladia_key is None and audio_url is None and audio_file is None:
        click.echo("Error: --audio-url or --audio-file is required.")
        return
    
    if not save_gladia_key:
        if gladia_key != "":

            if direct_translate and direct_translate_language is None:
                click.echo("Error: --direct-translate-language is required when using --direct-translate.")
                return 

            if diarization and diarization_max_speakers is None:
                click.echo("Error: --diarization-max-speakers is required when using --diarization.")
                return

            if audio_url is None and audio_file is None:
                click.echo("Error: --audio-url or --audio-file is required.")
                return
            else:
                click.echo("Transcribing audio file...")
                click.echo("This may take a few seconds, please wait...")
                headers = {
                    "accept": "application/json",
                    "x-gladia-key": gladia_key,
                }

                if output_format == "table":
                    this_output_format = "json"
                else:
                    this_output_format = output_format

                files = {
                    "language_behaviour": (None, language_behaviour),
                    "language": (None, language),
                    "toggle_noise_reduction": (None, "true" if noise_reduction else "false"),
                    "toggle_diarization": (None, "true" if diarization else "false"),
                    "diarization_max_speakers": (None, str(diarization_max_speakers)),
                    "toggle_direct_translate": (None, "true" if direct_translate else "false"),
                    "target_translation_language": (None, direct_translate_language),
                    "toggle_text_emotion_recognition": (None, "true" if text_emotion else "false"),
                    "toggle_summarization": (None, "true" if summarization else "false"),
                    "output_format": (None, this_output_format),
                }

                if audio_url:
                    files["audio_url"] = (None, audio_url)
                else:
                    mime = magic.Magic(mime=True)
                    file_type = mime.from_file(audio_file)
                    files["audio"] = (audio_file, open(audio_file, "rb"), file_type)

                response = requests.post(GLADIA_API_URL, headers=headers, files=files)

                if response.status_code != 200:
                    click.echo(f"Error: {response.status_code} - {response.text}")
                    return

                click.echo(f"{Color.BOLD}Transcript{Color.END}\n")

                if output_format == "table":
                    table = PrettyTable()
                    table.align = "l"
                    table.padding_width = 1
                    table.border = False
                    field_names = ["time_begin", "time_end", "confidence", "language"]

                    if diarization:
                        field_names.append("speaker")

                    if text_emotion:
                        field_names.append("emotion")

                    field_names.append("transcription")

                    table.field_names = field_names

                    for sentence in response.json()["prediction"]:
                        confidences = []
                        for words in sentence["words"]:
                            confidences.append(float(words["confidence"]))

                        # calculate the average
                        
                        confidence = round(sum(confidences) / len(confidences), 2)
                        


                        row = [
                            Color.GREEN + str("{:.3f}".format(sentence['time_begin'])) + Color.END, 
                            Color.GREEN + str("{:.3f}".format(sentence['time_end'])) + Color.END, 
                            Color.BLUE + str("{:.2f}".format(confidence)) + Color.END, 
                            Color.CYAN + sentence['language'] + Color.END, 
                        ]

                        if diarization:
                            row.append(Color.YELLOW + sentence['speaker'] + Color.END)

                        if text_emotion:
                            row.append(Color.YELLOW + sentence['emotion'] + Color.END)
                        
                        row.append(sentence['transcription'])
                        table.add_row(row)

                    click.echo(table)

                    if summarization:
                        click.echo("")
                        click.echo("=======")
                        click.echo("Summary")
                        click.echo("=======")
                        click.echo(response.json()["prediction_raw"]["summarization"])

                elif output_format == "json":
                    click.echo(response.json())
                else:
                    click.echo(response.json()["prediction"])
        else:
            click.echo("Error: Gladia API key not found.")
            click.echo("Please provide your Gladia API key using --gladia-key or save it using --save-gladia-key.")
            return

if __name__ == "__main__":
    transcribe()

