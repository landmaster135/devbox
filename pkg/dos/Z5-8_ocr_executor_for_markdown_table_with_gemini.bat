@echo off

set /p path="Input file or directory path for Markdown table generation with OCR: "
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

set /p token="Input token to request Gemini API: "

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\ocr-executor-with-ai.exe -path "%path%" -ai-type gemini -api-key %token% -generates-markdown-table
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
