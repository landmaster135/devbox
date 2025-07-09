@echo off

choice /c yn /n /m "Select whether to move original image files to archive directory or not  [y]='--move'  [n]='' : "
if %errorlevel% == 1 (
  set "moves=-move"
) else (
  set "moves="
)
echo %moves%

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\image-filterer.exe -src . -suffix blurred -x1 1690 -y1 1055 -x2 1872 -y2 1080 %moves% -mode blur -radius 50
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
