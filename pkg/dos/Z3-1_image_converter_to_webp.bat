@echo off
echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\image-converter.exe -src . -ext webp -q 80
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
