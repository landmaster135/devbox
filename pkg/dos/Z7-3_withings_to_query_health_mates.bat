@echo off
set /p start_date="Input start date to output (yyyy-MM-dd): "
set /p end_date="Input end date to output (yyyy-MM-dd): "
set /p client_id="Input your Client ID on Withings: "
set /p client_sec="Input your Client Secret on Withings: "
set /p user_id="Input your User ID on Withings: "
set /p access_token="Input your Access Token on Withings: "
set /p refresh_token="Input your Refresh Token on Withings: "

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\withings.exe -operation daily-summary -start-date %start_date% -end-date %end_date% -measure-types all -user-id %user_id% -client-id %client_id% -client-secret %client_sec% -access-token %access_token% -refresh-token %refresh_token% -output-file-path .\9_json_for_health_mates\%end_date%_out.json
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
