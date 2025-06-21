@echo off
set /p path="Input path of PDF file you wanna extract some files as JPG: "
set /p start="Input start page number (optional, press Enter to skip): "
set /p end="Input end page number (optional, press Enter to skip): "

echo --- プログラムを実行します ---

REM オプション文字列を構築
set options=
if not "%start%"=="" set options=%options% -start %start%
if not "%end%"=="" set options=%options% -end %end%

REM プログラムを実行
.\pkg\bin\win_amd64\pdf-merger.exe -extract %path% -output-dir .\888_images_from_pdf %options%

echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
