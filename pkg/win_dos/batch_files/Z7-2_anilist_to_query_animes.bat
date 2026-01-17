@echo off
set /p username="Input your AniList username: "
echo --- プログラムを実行します ---
.\pkg\bin\cli\win_amd64\anilist.exe -operation query-anime -username %username% -output-dir .
echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
