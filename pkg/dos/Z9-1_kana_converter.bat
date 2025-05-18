@echo off

set /p input_chr="Input characters you wanna convert: "
choice /c fh /n /m "Select conversion mode  [f]='full'  [h]='half' : "
if %errorlevel% == 1 (
  set "conv_mode=full"
) else (
  set "conv_mode=half"
)

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\kana-converter.exe -input %input_chr% -mode %conv_mode%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
