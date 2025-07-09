@echo off

set /p path="Input directory you wanna modify EXIF (e.g.: '.'): "
if /i "%path%"=="" (
  set "path=."
)

REM パスが存在するか確認
if not exist "%path%" (
  echo [ERROR] The specified path "%path%" does not exist.
  echo Input any key to exit...
  pause > nul
  exit /b 1
)

set /p ext="Select extension of files you wanna view  [p]='png' [j]='jpg,jpeg' [w]='webp' : "
if /i "%ext%"=="p" (
  set "ext=png"
) else if /i "%ext%"=="j" (
  set "ext=jpg,jpeg"
) else if /i "%ext%"=="w" (
  set "ext=webp"
) else (
  echo "Invalid choice. Input any key to exit..."
  pause > nul
  exit /b 1
)
echo "Selected extension: %ext%"

set /p method="Select option to modify EXIF values  [f]='--from-filename' [s]='--from-screenshot' : "
if /i "%method%"=="f" (
  set "method=--from-filename"
) else if /i "%method%"=="s" (
  set "method=--from-screenshot"
) else (
  echo "Invalid choice. Input any key to exit..."
  pause > nul
  exit /b 1
)
echo "Selected option: %method%"

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\exif-modifier.exe --dir %path% --ext %ext% %method% --verbose
.\pkg\bin\win_amd64\exif-viewer.exe -dir %path% -ext %ext% -props "File Modification Date/Time" -v
.\pkg\bin\win_amd64\exif-viewer.exe -dir %path% -ext %ext% -list-props
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
