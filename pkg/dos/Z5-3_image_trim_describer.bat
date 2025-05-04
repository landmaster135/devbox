@echo off

@REM set /p coordinates="Input 4 coordinates linked on each parameter to trim images (lefter-x, upper-y, righter-x, lower-y): "
@REM choice /c yn /n /m "Select whether to move original image files to archive directory or not  [y]='--move'  [n]='' : "
@REM if %errorlevel% == 1 (
@REM   set "moves=-move"
@REM ) else (
@REM   set "moves="
@REM )
@REM echo %moves%

echo --- プログラムを実行します ---
@REM .\pkg\bin\win_amd64\image_trim_describer.exe
.\image_trim_describer.exe
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
