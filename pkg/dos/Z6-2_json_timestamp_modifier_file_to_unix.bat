@echo off

set /p filepath="Input path of JSON file you wanna edit: "

choice /c ui /n /m "Select conversion format  [u]=--to-unix  [i]=--to-iso : "
rem ---- choice の戻り値 (errorlevel) で分岐 ----
rem  ※ /c ui だと i を押したとき errorlevel=2, u を押したとき errorlevel=1
if %errorlevel% == 2 (
  set "to=to-iso"
) else (
  set "to=to-unix"
)

set /p key="Input key to convert to ISO-8601 format or UNIX timestamp: "

choice /c yn /n /m "Select whether to enable the --is-jst option or not  [y]='--is-jst'  [n]='' : "
if %errorlevel% == 1 (
  set "is_jst=--is-jst"
) else (
  set "is_jst="
)

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\json-timestamp-modifier.exe %is_jst% %is_recursive% -file %filepath% -key %key% -mode %to%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
