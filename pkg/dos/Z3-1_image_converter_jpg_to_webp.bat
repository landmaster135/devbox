@echo off

choice /c yn /n /m "Select whether to move original image files to archive directory or not  [y]='--move'  [n]='' : "
if %errorlevel% == 1 (
  set "moves=-move"
) else (
  set "moves="
)

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\image-converter.exe -src . -ext webp -q 5 -archive .\5_original_files %moves%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
