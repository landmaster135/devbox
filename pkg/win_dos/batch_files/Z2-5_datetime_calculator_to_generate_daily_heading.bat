@echo off

setlocal enabledelayedexpansion
set /p offset="Input offset days in int type (negative int applicable): "

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\datetime-calculator.exe -operation generate-daily-heading -day-offset %offset%
endlocal

echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
