@echo off
set /p client_id="Input your Client ID on Withings: "
set /p client_sec="Input your Client Secret on Withings: "
set /p auth_code="Input your authorization code for Withings: "

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\withings.exe -operation request-token -client-id "%client_id%" -client-secret "%client_sec%" -authorization-code "%auth_code%" -redirect-uri "http://localhost:80/"
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
