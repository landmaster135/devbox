@echo off

set /p path="Input directory you wanna view EXIF: "
set /p ext="Select extension of files you wanna view  [p]='png' [j]='jpg,jpeg' [w]='webp' [m]='mp4,webm' : "
if /i "%ext%"=="p" (
  set "ext=png"
) else if /i "%ext%"=="j" (
  set "ext=jpg,jpeg"
) else if /i "%ext%"=="w" (
  set "ext=webp"
) else if /i "%ext%"=="m" (
  set "ext=mp4,webm"
) else (
  echo "Invalid choice. Input any key to exit..."
  pause > nul
  exit /b 1
)
echo "Selected extension: %ext%"

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\exif-viewer.exe -dir %path% -ext %ext% -max 4 -v
.\pkg\bin\win_amd64\exif-viewer.exe -dir %path% -ext %ext% -list-props
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
