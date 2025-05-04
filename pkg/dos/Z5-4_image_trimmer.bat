@echo off

set /p coordinates="Input 4 coordinates linked on each parameter to trim images (lefter-x, upper-y, righter-x, lower-y): "
choice /c yn /n /m "Select whether to move original image files to archive directory or not  [y]='--move'  [n]='' : "
if %errorlevel% == 1 (
  set "moves=-move"
) else (
  set "moves="
)
echo %moves%

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\image-trimmer.exe -src . -suffix cropped -x1 0 -y1 131 -x2 1080 -y2 2293 %moves%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
