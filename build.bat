@echo off

setlocal enabledelayedexpansion

set "cmd_opt="
for /F %%G in (requirements.txt) do (
  set "cmd_opt=!cmd_opt! --hidden-import=%%G"
)

echo py -m PyInstaller --onefile gladia_cli.py !cmd_opt!
py -m PyInstaller --onefile gladia_cli.py !cmd_opt!

