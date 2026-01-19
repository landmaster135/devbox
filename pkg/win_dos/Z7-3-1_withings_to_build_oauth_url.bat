@echo off
set /p client_id="Input your Client ID on Withings: "

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\withings.exe -operation auth-url -client-id "%client_id%" -redirect-uri "http://localhost:80/" -scope "user.metrics,user.activity" -state "random_state_value"
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
