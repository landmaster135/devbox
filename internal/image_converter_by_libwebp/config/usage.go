package config

const usageTemplate = `使用方法: %[1]s [オプション]

libwebp の cwebp を使って複数画像を WebP へ変換します。

オプション:
  --src-dir string      入力画像を探索するディレクトリ (デフォルト: .)
  --out-dir string      変換後画像の出力ディレクトリ (デフォルト: ./999_converted_images)
  --archive-dir string  処理済み元ファイルの退避先ディレクトリ (空なら退避しない)
  --move                退避時にコピーではなく移動する
  --ext string          出力形式 (デフォルト: webp、対応: webp)
  --q int               WebP 品質 1-100 (デフォルト: 99)
  --m int               cwebp 圧縮メソッド 0-6 (デフォルト: 4)
  --workers int         並列ワーカー数 (デフォルト: CPU数)
  --recursive, -R       サブディレクトリを再帰的に走査する
  --lossless            cwebp の lossless 圧縮を有効にする
  --help, -h            ヘルプを表示する

cwebp 固定オプション:
  -preset photo -metadata icc -sharp_yuv -progress -short
  --m は cwebp の -m として渡します。

使用例:
  %[1]s --src-dir ./photos --out-dir ./webp --q 99
  %[1]s --src-dir ./photos --out-dir ./webp --q 99 --m 6
  %[1]s --src-dir ./photos --out-dir ./webp --archive-dir ./archive --move --recursive

注意:
  cwebp が PATH に存在しない場合はエラー終了します。
`
