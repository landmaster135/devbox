@echo off
set url=%DB_SERVER_URL%
set endpoint="%url%/login"
set json=".\9_json\config\login_from_windows_01.json"
echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\http-request.exe -url %endpoint% -method "POST" -json %json%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
