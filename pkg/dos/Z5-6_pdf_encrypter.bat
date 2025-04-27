@echo off
set /p filepath="Input path of PDF file you wanna encrypt: "
set /p pw="Input password to encrypt PDF: "
echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\pdf-encrypter.exe -in %filepath% -mode encrypt -upw %pw% -opw %pw%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
