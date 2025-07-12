@echo off

choice /c vwa /n /m "Select how to rename image files  [v]='--vlc'  [w]='--win'  [a]='--android (screen_record)' : "
if %errorlevel% == 1 (
  set "method=-vlc"
) else if %errorlevel% == 2 (
  set "method=-win"
) else (
  set "method=-android"
)
echo %method%

choice /c yn /n /m "Select whether to move original image files to archive directory or not  [y]='--move'  [n]='' : "
if %errorlevel% == 1 (
  set "moves=-move"
) else (
  set "moves="
)

echo --- リネームプログラムを実行します ---
.\pkg\bin\cli\win_amd64\image-renamer-for-screenshot.exe -src . %method%
echo.
echo --- リネームプログラムの実行が完了しました ---

echo --- Webp変換プログラムを実行します ---
.\pkg\bin\cli\win_amd64\image-converter.exe -src . -ext webp -q 80 -archive .\5_original_files %moves%
echo.
echo --- Webp変換プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
