@echo off

setlocal enabledelayedexpansion

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\file-maneuver.exe --src-dirs \\wsl.localhost\Ubuntu-24.04\home\nov\devbox\pkg\bin\win_amd64 --dest-dir .\pkg\bin\win_amd64 --extensions "exe" --copy --overwrite
endlocal

echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
