@echo off

setlocal enabledelayedexpansion
set /p start="Input start of number for content of wine (WIxxxx_01.webp): "
choice /c nd /n /m "Select how to sort image files to rename  [n]='--name'  [t]='--time' : "
if %errorlevel% == 1 (
  set "sort=-name"
) else (
  set "sort=-time"
)
echo %sort%

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\image-renamer-for-content.exe -src . -operation "wine" -delimiter "" -suffix "" -start %start% %sort%
endlocal

echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
