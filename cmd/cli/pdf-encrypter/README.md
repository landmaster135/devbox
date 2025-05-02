# PDF Encrypter ツール使用方法

## 概要

PDF Encrypterは、PDFファイルの暗号化と復号化を行うためのコマンドラインツールです。AES-256暗号化を使用して、PDFファイルにパスワード保護を適用したり、解除したりすることができます。

## インストール

ビルドスクリプトを使用してツールをビルドします：

```bash
cd devbox
./scripts/build_pdf_encrypter.sh
```

ビルドが成功すると、以下の場所に実行ファイルが生成されます：

- Linux用: `pkg/bin/linux_amd64/pdf-encrypter`
- Windows用: `pkg/bin/win_amd64/pdf-encrypter.exe`

## 基本的な使い方

### PDFファイルの暗号化

PDFファイルを暗号化するには、以下のコマンドを使用します：

```bash
./pdf-encrypter -mode encrypt -in <入力PDFファイル> -out <出力PDFファイル> -upw <ユーザーパスワード> -opw <オーナーパスワード>
```

例：

```bash
./pdf-encrypter -mode encrypt -in plain.pdf -out locked.pdf -upw reader -opw owner
```

このコマンドは、`plain.pdf`を暗号化し、`locked.pdf`として保存します。ユーザーパスワード「reader」で閲覧でき、オーナーパスワード「owner」で管理操作（印刷、編集など）が可能になります。

- 出力ファイル名（-out）が指定されていない場合は、入力ファイルが上書きされます。
- オーナーパスワード（-opw）は必須です。
- ユーザーパスワード（-upw）は省略可能ですが、セキュリティのために設定することをお勧めします。

### PDFファイルの復号化

暗号化されたPDFファイルを復号化するには、以下のコマンドを使用します：

```bash
./pdf-encrypter -mode decrypt -in <暗号化されたPDFファイル> -out <出力PDFファイル> -opw <パスワード>
```

例：

```bash
./pdf-encrypter -mode decrypt -in locked.pdf -out plain.pdf -opw owner
```

このコマンドは、`locked.pdf`を復号化し、`plain.pdf`として保存します。パスワードには、ユーザーパスワードまたはオーナーパスワードを使用できます。

- 出力ファイル名（-out）が指定されていない場合は、入力ファイルが上書きされます。

## オプション

| オプション | 説明 |
|------------|------|
| `-mode` | 操作モード: `encrypt`（暗号化）または `decrypt`（復号化）（デフォルト: `encrypt`） |
| `-in` | 入力PDFファイルのパス（必須） |
| `-out` | 出力PDFファイルのパス（省略可能、省略時は入力ファイルを上書き） |
| `-upw` | ユーザーパスワード（閲覧用、暗号化時のみ使用） |
| `-opw` | オーナーパスワード（管理用、暗号化時は必須） |

## エラーメッセージ

以下のような場合にエラーメッセージが表示されます：

- 入力ファイルが指定されていない場合
- 暗号化モードでオーナーパスワードが指定されていない場合
- 指定されたファイルが存在しない場合
- 指定されたファイルがPDF形式でない場合
- パスワードが間違っている場合
- ファイルの読み書きに失敗した場合

## 使用例

### PDFファイルの暗号化（ユーザーパスワードとオーナーパスワードを設定）

```bash
./pdf-encrypter -mode encrypt -in document.pdf -out secure_document.pdf -upw user123 -opw admin456
```

### PDFファイルの暗号化（オーナーパスワードのみ設定）

```bash
./pdf-encrypter -mode encrypt -in document.pdf -opw admin456
```

### 暗号化されたPDFファイルの復号化

```bash
./pdf-encrypter -mode decrypt -in secure_document.pdf -out decrypted_document.pdf -opw admin456
```

## 注意事項

- 暗号化には AES-256 暗号化アルゴリズムが使用されます。
- パスワードを忘れた場合、ファイルを復号化することはできません。
- 出力PDFファイルは、既存のファイルがある場合は上書きされます。
- 暗号化されたPDFファイルは、Adobe Acrobat Reader などの標準的なPDFビューアで開くことができます。
- 暗号化の強度は設定したパスワードの複雑さに依存します。強力なパスワードを使用することをお勧めします。
