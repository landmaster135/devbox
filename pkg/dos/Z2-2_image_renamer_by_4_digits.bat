@echo off

setlocal enabledelayedexpansion
set /p prefix="Input prefix to rename image files (e.g. yyyymmdd_0000.jpg): "
set /p start="Input number to renumber the images  (e.g. !prefix!_xxxx.jpg): "
choice /c nd /n /m "Select how to sort image files to rename  [n]='--name'  [d]='--time' : "
if %errorlevel% == 1 (
  set "sort=-name"
) else (
  set "sort=-time"
)
echo %sort%

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\image-renamer.exe -src . -digits 4 -delimiter "_" -prefix %prefix% -start %start% %sort%
endlocal

echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
