@echo off

echo --- プログラムを実行します ---
.\pkg\bin\win_amd64\movie-converter-for-webm.exe -input-dir . -input-ext mp4 -output-dir .\999_converted_images -output-ext webm -conversion-mode cbr -audio-codec opus -video-quality 35
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
