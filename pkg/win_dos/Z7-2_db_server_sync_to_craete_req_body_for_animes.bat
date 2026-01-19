@echo off
set /p json="Input your json file path: "
set /p out_path="Input output file path: "
echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\db-server-sync.exe -operation append-anime -input-file-path %json% -output-file-path %out_path%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
