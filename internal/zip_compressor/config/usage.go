package config

import (
	flagParser "github.com/landmaster135/devbox/internal/zip_compressor/infrastructures/flag_parser"
)

const usageTemplate = `Zip圧縮CLIツール

使用方法:
  ファイル/ディレクトリ圧縮:
    %[1]s -operation compress -path /path/to/file_or_directory
    %[1]s -o compress -p /path/to/file_or_directory
    %[1]s compress /path/to/file_or_directory

  Zipファイル展開:
    %[1]s -operation decompress -path /path/to/archive.zip
    %[1]s -o decompress -p /path/to/archive.zip
    %[1]s decompress /path/to/archive.zip

オプション:
  -operation, -o    操作タイプ (compress, decompress)
  -path, -p         対象ファイル/ディレクトリのパス
  -help, -h         このヘルプを表示

例:
  # ファイル圧縮
  %[1]s compress /home/user/document.txt
  # -> document.txt.zip が作成される

  # ディレクトリ圧縮
  %[1]s compress /home/user/my_folder
  # -> my_folder.zip が作成される

  # Zipファイル展開
  %[1]s decompress /home/user/archive.zip
  # -> archive_decompressed/ ディレクトリに展開される

`

// PrintUsage は使用方法を表示する
func PrintUsage() {
	flagParser.PrintUsage(usageTemplate)
}
