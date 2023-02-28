from setuptools import setup

setup(
    name='gladia-transcriber',
    version='0.1.0',
    author='Jean-Louis Queguiner',
    author_email='jlqueguiner@gladia.io',
    description='Transcribe audio files using the Gladia API',
    packages=['gladia_transcriber', 'gladia_transcriber.config'],
    package_data={'gladia_transcriber.config': ['config.ini']},
    install_requires=[
        'click',
        'requests',
        'prettytable',
    ],
    entry_points={
        'console_scripts': [
            'gladia-transcriber = gladia_transcriber.transcribe:transcribe',
        ],
    },
)
