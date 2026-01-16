#!/bin/bash

# Select rename operation
echo -n "Select how to rename image files  [v]='vlc'  [w]='win'  [x]='xiaomi'  [a]='pixel (screen_record)' : "
read -n 1 choice
echo

case $choice in
  v|V)
    method="-operation=vlc"
    ;;
  w|W)
    method="-operation=win"
    ;;
  x|X)
    method="-operation=xiaomi"
    ;;
  a|A)
    method="-operation=pixel"
    ;;
  *)
    echo "Invalid choice. Exiting."
    exit 1
    ;;
esac
echo "$method"

# Select whether to move files
echo -n "Select whether to move original image files to archive directory or not  [y]='--move'  [n]='' : "
read -n 1 move_choice
echo

if [[ $move_choice == "y" || $move_choice == "Y" ]]; then
  moves="-move"
else
  moves=""
fi

echo "--- リネームプログラムを実行します ---"
./pkg/bin/cli/linux_amd64/image-renamer-for-screenshot -src . "$method"
echo
echo "--- リネームプログラムの実行が完了しました ---"

echo "--- Webp変換プログラムを実行します ---"
./pkg/bin/cli/linux_amd64/image-converter -src . -ext webp -q 80 -archive ./5_original_files $moves
echo
echo "--- Webp変換プログラムの実行が完了しました ---"
echo "--- 何かキーを押すと終了します ---"
read -n 1 -s
