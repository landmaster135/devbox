@echo off

setlocal enabledelayedexpansion
set /p prefix="Input prefix of pictures for LINE (yyyyMMddHH): "
choice /c nd /n /m "Select how to sort image files to rename  [n]='--name'  [d]='--time' : "
if %errorlevel% == 1 (
  set "sort=-name"
) else (
  set "sort=-time"
)
echo %sort%

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\image-renamer.exe -src .\1-3_image_renamer_with_prefix_input -digits 4 -delimiter "" -prefix %prefix% -start 1 %sort%
endlocal

echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
