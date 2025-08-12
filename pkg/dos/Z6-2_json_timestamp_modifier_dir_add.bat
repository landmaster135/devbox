@echo off
set /p dirpath="Input dir path having JSON files you wanna edit: "
set /p key="Input key to add UNIX timestamp: "

choice /c yn /n /m "Select whether to enable the --recursive option or not  [y]='--recursive'  [n]='' : "
if %errorlevel% == 1 (
  set "is_recursive=-recursive"
) else (
  set "is_recursive="
)

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\json-timestamp-modifier.exe -dir %dirpath% %is_recursive% -key %key% -mode add
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
