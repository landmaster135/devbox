@echo off
set /p client_id="Input your Client ID on Withings: "
set /p client_sec="Input your Client Secret on Withings: "
set /p refresh_token="Input your Refresh Token on Withings: "

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\withings.exe -operation refresh-token -client-id "%client_id%" -client-secret %client_sec% -refresh-token "%refresh_token%"
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
