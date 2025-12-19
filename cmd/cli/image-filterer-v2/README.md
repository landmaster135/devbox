# Image Filterer V2

CPUベースで画像フィルタ処理を行うCLIツールです。グレースケール化、ティント付与、ビネット効果といった基本的なエフェクトをワンショットで適用できます。

## 概要

指定した入力画像を読み込み、選択したフィルタと強度に応じて加工した結果をファイルへ書き出します。処理は標準ライブラリのみで完結します。PNG/JPEGの入力・出力に対応し、出力先が未指定の場合は自動的にファイル名を補完します。

## 機能

- **グレースケール変換**: 輝度の加重平均を用いた自然なモノクロ化
- **カラー調整**: HEXカラー指定によるティント付与と強度調整
- **ビネット効果**: 画像周辺を滑らかに減光させるエフェクト
- **出力制御**: 既存ファイルを書き換える／別名で保存する／フォーマットを切り替える
- **自動パラメータ補完**: 未指定の出力先やティントカラーを安全な既定値で補完

## インストール

```bash
# プロジェクトルートで実行
go build -o bin/image-filterer-v2 ./cmd/cli/image-filterer-v2
```

## 使用方法

### 基本的な使用方法

```bash
# 直接実行
./bin/image-filterer-v2 -input ./sample_data/input.png -mode grayscale -strength 0.8

# go run で実行
go run ./cmd/cli/image-filterer-v2 -input ./sample_data/input.png -mode colorize -tint "#ff8844" -strength 0.6
```

### オプション

| オプション | 説明 | 必須 | 既定値 | 例 |
|------------|------|------|--------|-----|
| `-input` | 入力画像ファイルのパス | * | なし | `-input ./image.png` |
| `-output` | 出力ファイルのパス | - | 入力名 + `_filtered` | `-output ./dist/result.png` |
| `-format` | 出力フォーマット（拡張子が無い場合に使用） | - | 入力拡張子または `png` | `-format jpg` |
| `-mode` | 適用するフィルタ（`grayscale`/`colorize`/`vignette`） | - | `grayscale` | `-mode vignette` |
| `-strength` | フィルタ強度（0.0〜1.0） | - | `1.0` | `-strength 0.7` |
| `-tint` | `colorize` モード用ティントカラー（HEX） | - | `#ffffff` | `-tint "#33aaff"` |

### フィルタごとの挙動

| モード | 説明 | 主な調整項目 |
|--------|------|--------------|
| `grayscale` | 輝度加重平均によるモノクロ化 | `-strength` で元画像とのブレンド比率を変更 |
| `colorize` | 指定カラーへのティント付与 | `-tint` で色、`-strength` で寄せ具合を調整 |
| `vignette` | 周辺減光エフェクト | `-strength` が大きいほど中央以外が暗くなる |

## 使用例

```bash
# グレースケールに変換し、結果を自動生成ファイルへ保存
./bin/image-filterer-v2 -input ./photos/city.png -mode grayscale

# ティントを適用し、JPEGで出力
go run ./cmd/cli/image-filterer-v2 \
  -input ./photos/sunrise.png \
  -output ./exports/sunrise-warm.jpg \
  -format jpg \
  -mode colorize \
  -tint "#ff7f50" \
  -strength 0.65

# ビネット効果を抑えめに適用
./bin/image-filterer-v2 -input ./photos/forest.png -mode vignette -strength 0.4
```

処理が完了すると、標準出力に保存先パスが表示されます。
