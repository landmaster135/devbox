@echo off

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\movie-converter-for-gif.exe -input-dir . -input-ext mp4 -output-dir .\999_converted_images -output-ext gif
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
