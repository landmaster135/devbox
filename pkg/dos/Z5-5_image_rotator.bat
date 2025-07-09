@echo off

choice /c abc /n /m "Select which to rotate image files to rotate  [a]='right-90'  [b]='left-90'  [c]='180' : "
if %errorlevel% == 1 (
  set "angle=-angle 90"
) else if %errorlevel% == 2 (
  set "angle=-angle 270"
) else (
  set "angle=-angle 180"
)
echo %moves%
choice /c yn /n /m "Select whether to move original image files to archive directory or not  [y]='--move'  [n]='' : "
if %errorlevel% == 1 (
  set "moves=-move"
) else (
  set "moves="
)
echo %moves%

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\image-rotator.exe -src . -suffix rotated %angle% %moves%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
