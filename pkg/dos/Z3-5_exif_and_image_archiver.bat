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
@REM if /i "%ext%"=="p" (
@REM   set "ext=png"
@REM ) else if /i "%ext%"=="j" (
@REM   set "ext=jpg,jpeg"
@REM ) else if /i "%ext%"=="w" (
@REM   set "ext=webp"
@REM ) else (
@REM   echo "Invalid choice. Input any key to exit..."
@REM   pause > nul
@REM   exit /b 1
@REM )
@REM echo "Selected extension: %ext%"

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\image-renamer-with-exif.exe --dir %path% --ext jpg --verbose
.\pkg\bin\win_amd64\image-converter.exe -src . -ext webp -q 80 -archive .\5_original_files %moves%


@REM .\pkg\bin\win_amd64\exif-modifier.exe --dir %path% --ext %ext% %method% --verbose
@REM .\pkg\bin\win_amd64\exif-viewer.exe -dir %path% -ext %ext% -props "File Modification Date/Time" -v
@REM .\pkg\bin\win_amd64\exif-viewer.exe -dir %path% -ext %ext% -list-props
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
