@echo off

set /p path="Input file or directory path for OCR: "
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

set /p prompt="Input choice of prompt what you generate [p] plain text or [t] table: "
if /i "%prompt%"=="p" (
  set "prompt=OCRして。"
) else if /i "%prompt%"=="t" (
  set "prompt=OCRして、Markdownのテーブル形式にして。"
) else (
  echo "Invalid choice. Input any key to exit..."
  pause > nul
  exit /b 1
)

set /p token="Input token to request Gemini API: "

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\ocr-executor-with-ai.exe -path "%path%" -ai-type gemini -api-key %token% -prompt "%prompt%"
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
