@echo off

set /p filepath="Input path of JSON file you wanna edit: "
set /p key="Input key to set: "
set /p value="Input value to set: "

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\json-modifier.exe -file %filepath% -key %key% -set %value%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
