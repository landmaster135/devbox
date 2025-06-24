@echo off

@REM set /p path="Input directory you wanna rename media files (e.g.: '.'): "
@REM if /i "%path%"=="" (
@REM   set "path=."
@REM )

set path_01_01=.\1-1_image_renamer_with_exif
set path_01_02=.\1-2_image_renamer_for_screenshot
set path_01_05=.\1-4_image_converter_to_webp
set path_01_06=.\1-4_image_converter_to_webp_org
set path_01_07=.\1-5_image_archiver_terminated

REM パスが存在するか確認
if not exist "%path_01_01%" (
  echo [ERROR] The specified path "%path_01_01%" does not exist.
  echo Input any key to exit...
  pause > nul
  exit /b 1
)
if not exist "%path_01_02%" (
  echo [ERROR] The specified path "%path_01_02%" does not exist.
  echo Input any key to exit...
  pause > nul
  exit /b 1
)
if not exist "%path_01_05%" (
  echo [ERROR] The specified path "%path_01_05%" does not exist.
  echo Input any key to exit...
  pause > nul
  exit /b 1
)
if not exist "%path_01_06%" (
  echo [ERROR] The specified path "%path_01_06%" does not exist.
  echo Input any key to exit...
  pause > nul
  exit /b 1
)
if not exist "%path_01_07%" (
  echo [ERROR] The specified path "%path_01_07%" does not exist.
  echo Input any key to exit...
  pause > nul
  exit /b 1
)


echo --- プログラムを実行します ---
echo "===============  EXIF modification: Part 1  ============================================"
@REM 任意の日時を文字列で入力して自動採番する処理
echo "===================================================================================="
echo "===============  JPG renaming  ====================================================="
.\pkg\bin\win_amd64\image-renamer-with-exif.exe --dir %path_01_01% --ext jpg --verbose
echo "===================================================================================="
echo "===============  MP4 renaming  ====================================================="
.\pkg\bin\win_amd64\image-renamer-with-exif.exe --dir %path_01_01% --ext mp4 --verbose
echo "===================================================================================="
echo "===============  Screenshot (PNG, MP4) renaming  ======================================"
.\pkg\bin\win_amd64\image-renamer-for-screenshot.exe -src %path_01_02% -to-datetime
echo "===================================================================================="
echo "===============  EXIF modification: Part 2  ============================================"
.\pkg\bin\win_amd64\exif-modifier.exe --dir %path_01_02% --from-filename --ext png --verbose
.\pkg\bin\win_amd64\exif-modifier.exe --dir %path_01_02% --from-filename --ext mp4 --verbose
echo "===================================================================================="
echo "===============  WEBP conversion  ====================================================="
.\pkg\bin\win_amd64\image-converter.exe -src %path_01_01% -out .\%path_01_05% -ext webp -q 70 -archive .\%path_01_05%_org -move
.\pkg\bin\win_amd64\image-converter.exe -src %path_01_02% -out .\%path_01_05% -ext webp -q 70 -archive .\%path_01_05%_org -move
echo "===================================================================================="
echo "===============  EXIF mirroring  ====================================================="
.\pkg\bin\win_amd64\exif-mirror.exe --source-dir %path_01_06% --target-dir %path_01_05% --source-ext jpg --target-ext webp
.\pkg\bin\win_amd64\exif-mirror.exe --source-dir %path_01_06% --target-dir %path_01_05% --source-ext png --target-ext webp
echo "===================================================================================="
echo "===============  EXIF viewer (WEBP)  ================================================="
.\pkg\bin\win_amd64\exif-viewer.exe -dir %path_01_05% -ext webp -props "File Modification Date/Time"
.\pkg\bin\win_amd64\exif-viewer.exe -dir %path_01_05% -ext webp -list-props
echo "===================================================================================="
echo "===============  EXIF viewer (MP4)  ================================================="
.\pkg\bin\win_amd64\exif-viewer.exe -dir %path_01_01% -ext mp4 -props "File Modification Date/Time"
.\pkg\bin\win_amd64\exif-viewer.exe -dir %path_01_01% -ext mp4 -list-props
echo "===================================================================================="
echo "===============  EXIF viewer (MP4 SS)  ================================================"
.\pkg\bin\win_amd64\exif-viewer.exe -dir %path_01_02% -ext mp4 -props "File Modification Date/Time"
.\pkg\bin\win_amd64\exif-viewer.exe -dir %path_01_02% -ext mp4 -list-props
echo "===================================================================================="
echo "===============  File moving  ======================================================"
.\pkg\bin\win_amd64\file-maneuver.exe --src-dirs %path_01_05% --extensions webp --dest-dir %path_01_07%
.\pkg\bin\win_amd64\file-maneuver.exe --src-dirs %path_01_01%,%path_01_02% --extensions mp4 --dest-dir %path_01_07%
echo "===================================================================================="

echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
