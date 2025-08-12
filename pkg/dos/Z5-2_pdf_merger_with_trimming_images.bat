@echo off

set /p purpose="Select how to trim images by purpose taking screenshots  [g]='General' [c]='Comic' : "
if /i "%purpose%"=="g" (
  set "coordinate_x1=0"
  set "coordinate_y1=131"
  set "coordinate_x2=1080"
  set "coordinate_y2=2293"
) else if /i "%purpose%"=="c" (
  set "coordinate_x1=0"
  set "coordinate_y1=341"
  set "coordinate_x2=1080"
  set "coordinate_y2=2056"
) else (
  echo "Invalid choice. Input any key to exit..."
  pause > nul
  exit /b 1
)
set /p quality="Input integer (75 as default, 30 useful for comic) for quality of output jpg (1-100): "

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\image-trimmer.exe -src . -suffix cropped -x1 %coordinate_x1% -y1 %coordinate_y1% -x2 %coordinate_x2% -y2 %coordinate_y2% -move
.\pkg\bin\cli\win_amd64\image-converter.exe -src . -ext jpg -q %quality% -archive .\5_original_files -move
.\pkg\bin\cli\win_amd64\pdf-merger.exe -dir %HOMEPATH%\Downloads\picture_backup
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
