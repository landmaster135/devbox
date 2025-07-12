@echo off

set coordinates="-x1 0 -y1 131 -x2 1080 -y2 2293"
choice /c a /n /m "Select how to trim images by devices taking screenshots  [a]='Android' : "
if %errorlevel% == 1 (
  set "coordinate_x1=0"
  set "coordinate_y1=131"
  set "coordinate_x2=1080"
  set "coordinate_y2=2293"
) else (
  set "coordinate_x1="
  set "coordinate_y1="
  set "coordinate_x2="
  set "coordinate_y2="
)
choice /c yn /n /m "Select whether to move original image files to archive directory or not  [y]='--move'  [n]='' : "
if %errorlevel% == 1 (
  set "moves=-move"
) else (
  set "moves="
)
echo %moves%

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\image-trimmer.exe -src . -suffix cropped -x1 %coordinate_x1% -y1 %coordinate_y1% -x2 %coordinate_x2% -y2 %coordinate_y2% %moves%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
