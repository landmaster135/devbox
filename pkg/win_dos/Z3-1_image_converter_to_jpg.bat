@echo off

set /p quality="Input integer (75 as default) for quality of output jpg (1-100): "
choice /c yn /n /m "Select whether to move original image files to archive directory or not  [y]='--move'  [n]='' : "
if %errorlevel% == 1 (
  set "moves=-move"
) else (
  set "moves="
)
echo %moves%

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\image-converter.exe -src . -ext jpg -q %quality% -archive .\5_original_files %moves%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
