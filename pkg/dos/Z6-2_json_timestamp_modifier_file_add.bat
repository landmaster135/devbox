@echo off
set /p filepath="Input path of JSON file you wanna edit: "
set /p key="Input key to add UNIX timestamp: "
echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\json-timestamp-modifier.exe -file %filepath% -key %key% -mode add
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
