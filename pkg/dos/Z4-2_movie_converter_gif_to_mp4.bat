@echo off

echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\movie-converter-for-gif.exe -input-dir . -input-ext gif -output-dir .\999_converted_images -output-ext mp4
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
