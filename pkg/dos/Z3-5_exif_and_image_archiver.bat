@echo off

@REM set /p path="Input directory you wanna rename media files (e.g.: '.'): "
@REM if /i "%path%"=="" (
@REM   set "path=."
@REM )

set path_01_01=.\1-1_image_renamer_with_exif
set path_01_02=.\1-1_image_renamer_for_screenshot
set path_01_05=.\1-2_image_converter_to_webp
set path_01_06=.\1-2_image_converter_to_webp_org

REM パスが存在するか確認
if not exist "%path_01_01%" (
  echo [ERROR] The specified path "%path_01_01%" does not exist.
  echo Input any key to exit...
  pause > nul
  exit /b 1
)
if not exist "%path_01_02%" (
  echo [ERROR] The specified path "%path_01_01%" does not exist.
  echo Input any key to exit...
  pause > nul
  exit /b 1
)

@REM set /p ext="Select extension of files you wanna view  [p]='png' [j]='jpg,jpeg' [w]='webp' : "
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
echo "===============  JPG renaming  ====================================================="
.\pkg\bin\win_amd64\image-renamer-with-exif.exe --dir %path_01_01% --ext jpg --verbose
echo "===================================================================================="
echo "===============  MP4 renaming  ====================================================="
.\pkg\bin\win_amd64\image-renamer-with-exif.exe --dir %path_01_01% --ext mp4 --verbose
echo "===================================================================================="
echo "===============  Screenshot (PNG, MP4) renaming  ======================================"
@REM .\pkg\bin\win_amd64\image-renamer-with-exif.exe --dir %path_01_01% --ext png --verbose
.\pkg\bin\win_amd64\image-renamer-for-screenshot.exe -src %path_01_02% -to-datetime
echo "===================================================================================="
echo "===============  EXIF modification  ==================================================="
.\pkg\bin\win_amd64\exif-modifier.exe --dir %path_01_02% --from-filename --ext png --verbose
.\pkg\bin\win_amd64\exif-modifier.exe --dir %path_01_02% --from-filename --ext mp4 --verbose
echo "===================================================================================="
echo "===============  WEBP conversion  ====================================================="
@REM .\pkg\bin\win_amd64\image-converter.exe -src %path_01_01% -out .\%path_01_05% -ext webp -q 80 -archive .\%path_01_05%_org -move
@REM .\pkg\bin\win_amd64\image-converter.exe -src %path_01_02% -out .\%path_01_05% -ext webp -q 80 -archive .\%path_01_05%_org -move
echo "===================================================================================="
echo "===============  EXIF mirroring  ====================================================="
@REM .\pkg\bin\win_amd64\exif-mirror.exe --source-dir %path_01_06% --target-dir %path_01_05% --source-ext jpg --target-ext webp
@REM .\pkg\bin\win_amd64\exif-mirror.exe --source-dir %path_01_06% --target-dir %path_01_05% --source-ext png --target-ext webp
echo "===================================================================================="

@REM .\pkg\bin\win_amd64\exif-modifier.exe --dir %path% --ext %ext% %method% --verbose
@REM .\pkg\bin\win_amd64\exif-viewer.exe -dir %path% -ext %ext% -props "File Modification Date/Time" -v
@REM .\pkg\bin\win_amd64\exif-viewer.exe -dir %path% -ext %ext% -list-props
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
