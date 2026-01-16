@echo off
set /p path="Input path of PDF file you wanna extract some files as JPG: "
echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\pdf-merger.exe -extract %path% -output-dir .\888_images_from_pdf
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
