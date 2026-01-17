@echo off
echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\pdf-merger.exe -dir %HOMEPATH%\Downloads\picture_backup
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
