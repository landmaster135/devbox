@echo off

set /p coordinates="Input 4 coordinates linked on each parameter to trim images (-x1, -y1, -x2, -y2): "
choice /c yn /n /m "Select whether to move original image files to archive directory or not  [y]='--move'  [n]='' : "
if %errorlevel% == 1 (
  set "moves=-move"
) else (
  set "moves="
)
echo %moves%

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\image-filterer.exe -src . -suffix blurred %coordinates% %moves% -mode blur -radius 50
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
