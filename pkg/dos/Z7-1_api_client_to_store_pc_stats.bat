@echo off
set url=%DB_SERVER_URL%
set endpoint="%url%/pc-stats/append"
set /p json="Input your json file path: "
set /p token="Input your token: "
echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\api-client.exe -url %endpoint% -method "POST" -json %json% -token %token%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
