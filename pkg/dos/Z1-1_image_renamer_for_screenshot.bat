@echo off

choice /c vw /n /m "Select how to rename image files  [v]='--vlc'  [w]='--win' : "
if %errorlevel% == 1 (
  set "method=-vlc"
) else (
  set "method=-win"
)
echo %method%

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\image-renamer-for-screenshot.exe -src . %method%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
