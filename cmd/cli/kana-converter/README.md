# カナ変換ツール (Kana Converter)

カタカナを含む文字列を全角または半角に変換するコマンドラインツールです。

## 機能

- 半角カタカナを全角カタカナに変換
- 全角カタカナを半角カタカナに変換
- コマンドライン引数からの入力サポート
- フラグまたは位置引数での入力指定

## 使い方

```bash
kana-converter [オプション] [入力文字列]
```

### オプション

| オプション | デフォルト値 | 説明 |
|------------|--------------|------|
| `-input` | `""` (空) | 変換する入力文字列 |
| `-mode` | `full` | 変換モード (`full` = 全角変換, `half` = 半角変換) |

### 使用例

#### 位置引数で半角カナを全角カナに変換（デフォルトのfullモード）

```bash
kana-converter ｶﾀｶﾅ
```

出力:
```
Input characters: ｶﾀｶﾅ
Mode: full-width
Output characters: カタカナ
```

#### フラグで半角カナを全角カナに変換

```bash
kana-converter -input="ｶﾀｶﾅ" -mode=full
```

出力:
```
Input characters: ｶﾀｶﾅ
Mode: full-width
Output characters: カタカナ
```

#### 全角カナを半角カナに変換

```bash
kana-converter -input="カタカナ" -mode=half
```

出力:
```
Input characters: カタカナ
Mode: half-width
Output characters: ｶﾀｶﾅ
```

#### スペースを含む文字列の変換

```bash
kana-converter "ｶﾀｶﾅ ﾃｽﾄ"
```

出力:
```
Input characters: ｶﾀｶﾅ ﾃｽﾄ
Mode: full-width
Output characters: カタカナ テスト
```

## 技術詳細

このツールは次のGolangパッケージを使用しています：

- `golang.org/x/text/unicode/norm` - Unicode正規化
- `golang.org/x/text/width` - 文字幅変換

## ビルド方法

リポジトリのルートディレクトリで以下のコマンドを実行します：

```bash
cd devbox
go build -o bin/kana-converter ./cmd/cli/kana-converter
```

## 依存関係のインストール

必要なパッケージをインストールするには：

```bash
go get -u golang.org/x/text
```

## アーキテクチャ

このプログラムは以下のコンポーネントで構成されています：

1. `main.go` - CLIインターフェースとエントリーポイント
2. `internal/independencies/kana_converter/usecases` - カナ変換のコアロジック

変換機能は再利用可能なパッケージとして実装されており、他のプログラムからも利用できます。
