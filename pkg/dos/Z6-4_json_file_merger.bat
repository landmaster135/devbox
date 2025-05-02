@echo off
set /p dirpath="Input directory path having JSON files you wanna edit: "
set /p key="Input key to store JSON objects merged with this process: "
set /p output="Input name of new file merged all JSON files (Don't forget '.json'): "
echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\json-file-merger.exe -dir %dirpath% -key %key% -output %output%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
