@echo off

setlocal enabledelayedexpansion
set /p path="Input .txt or .md file path to extract working time: "
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
.\pkg\bin\cli\win_amd64\datetime-calculator.exe -operation parse-time -file-path %path% -output-unit minute
endlocal

echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
