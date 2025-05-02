@echo off
set /p filepath="Input path of PDF file you wanna decrypt: "
set /p pw="Input password to decrypt PDF: "
echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\pdf-encrypter.exe -in %filepath% -mode decrypt -opw %pw%
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
