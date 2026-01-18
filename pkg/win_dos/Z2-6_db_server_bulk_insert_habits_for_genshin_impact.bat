@echo off

setlocal enabledelayedexpansion
set entity="habits"
set process="bulk-insert-by-category"
set category="playing_genshin_impact"
set list=".\0000_tmp_draft.txt"
set output=".\0000_tmp_append_request.json"

echo --- プログラムを実行します ---
.\pkg\bin\db-server\cli\win_amd64\all_cli_entry.exe --entity %entity% --process %process% -- --category %category% --list %list% --output %output%
endlocal

echo.
echo --- プログラムの実行が完了しました ---
echo --- 何かキーを押すと終了します ---
pause > nul
