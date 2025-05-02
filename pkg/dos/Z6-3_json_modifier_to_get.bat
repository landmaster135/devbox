@echo off

set /p filepath="Input path of JSON file you wanna edit: "
set /p key="Input key to get: "
choice /c yn /n /m "Select whether to get all keys and values or not  [y]='-get-all'  [n]='-get' : "
if %errorlevel% == 1 (
  set "gets_all=-get-all"
) else (
  set "gets_all=-get"
)

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\json-modifier.exe -file %filepath% -key %key% %gets_all%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
