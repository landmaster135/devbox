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

set "prompt=この画像からキャラクターのステータスを読み取って、Markdownのテーブル形式にして。補足や説明は不要です。"
set "system_instruction=以下のカラムを持ったテーブル形式で出力して。\n---\nCharacter Name\nLevel\nHeart\nAttack\nDefense\nWork Speed\nPassive Skill 1\nPassive Skill 2\nPassive Skill 3\nPassive Skill 4\nPotential of Heart\nPotential of Attack"

set /p token="Input token to request Gemini API: "

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\ocr-executor-with-ai.exe -path "%path%" -ai-type gemini -api-key %token% -prompt "%prompt%" -system-instruction %system_instruction%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
