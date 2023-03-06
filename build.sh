#!/bin/bash

cmd_opt=""
while read package; do
  cmd_opt+=" --hidden-import=$package"
done < requirements.txt

pyinstaller --onefile gladia_cli.py $cmd_opt 

