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

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\image-renamer-for-screenshot.exe -src . %method%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
