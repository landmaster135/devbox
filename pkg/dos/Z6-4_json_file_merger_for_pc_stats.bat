@echo off
set dirpath=".\9_json"
set key="pc_stats"
set output=".\9_json\__merged_%date:~,4%%date:~5,2%%date:~8,2%.json"
echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\json-file-merger.exe -dir %dirpath% -key %key% -output %output%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
