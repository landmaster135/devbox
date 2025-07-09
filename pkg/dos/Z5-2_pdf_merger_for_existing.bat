@echo off

set /p path="Input PDF file you wanna add new pages (e.g.: '.'): "
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

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\pdf-merger.exe -dir %HOMEPATH%\Downloads\picture_backup -add %path%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
