@echo off

setlocal enabledelayedexpansion
set /p src_dir="Input src-dir: "
set /p dest_dir="Input dest-dir: "

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\file-maneuver.exe --src-dirs %src_dir% --dest-dir %dest_dir% --extensions "webp" --name-contains "thumbnail_of_game"
endlocal

echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
