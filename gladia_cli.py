import click
import requests
import os
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
@click.option("--language", default="english", help="Language spoken in the audio file.")
@click.option("--language-behaviour", default="automatic multiple languages", help="Determines how to handle multi-language audio.")
@click.option("--noise-reduction", is_flag=True, help="Apply noise reduction to the audio.")
@click.option("--output-format", default="json", help="Format in which to return the transcription results.")
@click.option("--diarization", is_flag=True, help="Perform speaker diarization.")
@click.option("--gladia-key", help="API key for Gladia. Get it at https://app.gladia.io/account")
@click.option("--save-gladia-key", is_flag=True, help="Save the API key to a configuration file.")
def transcribe(audio_url, language, language_behaviour, noise_reduction, output_format, diarization, gladia_key, save_gladia_key):
    """
    Transcribe an audio file using the Gladia API.
    """
    if gladia_key is None:
        gladia_key = get_gladia_key()

    if save_gladia_key is True:
        save_gladia_key_to_file(gladia_key)
    
    if gladia_key is None and not save_gladia_key:
        click.echo("Error: Gladia API key not found.")
        return
    
    if save_gladia_key is None and audio_url is None:
        click.echo("Error: --audio-url is required.")
        return
    
    if not save_gladia_key:
        if gladia_key != "":
            if audio_url is None:
                click.echo("Error: --audio-url is required.")
                return
            else:
                click.echo("Transcribing audio file...")
                headers = {
                    "accept": "application/json",
                    "x-gladia-key": gladia_key,
                }

                files = {
                    "language": (None, language),
                    "language_behaviour": (None, language_behaviour),
                    "noise_reduction": (None, "true" if noise_reduction else "false"),
                    "output_format": (None, output_format),
                    "toogle_diarization": (None, "true" if diarization else "false"),
                }

                if audio_url:
                    files["audio_url"] = (None, audio_url)

                response = requests.post(GLADIA_API_URL, headers=headers, files=files)

                if response.status_code != 200:
                    click.echo(f"Error: {response.status_code} - {response.text}")
                    return

                click.echo(f"{Color.BOLD}Transcript{Color.END}\n")

                table = PrettyTable()
                table.field_names = ["time_begin", "time_end", "probability", "language", "speaker", "transcription"]

                table.align = "l"
                table.padding_width = 1
                table.border = False

                for sentence in response.json()["prediction"]:
                    table.add_row([
                        Color.GREEN + str(sentence['time_begin']) + Color.END, 
                        Color.GREEN + str(sentence['time_end']) + Color.END, 
                        Color.BLUE + str(sentence['probability']) + Color.END, 
                        Color.CYAN + sentence['language'] + Color.END, 
                        sentence['speaker'], 
                        sentence['transcription']
                        ])
                click.echo(table)
        else:
            click.echo("Error: Gladia API key not found.")
            click.echo("Please provide your Gladia API key using --gladia-key or save it using --save-gladia-key.")
            return

if __name__ == "__main__":
    transcribe()

