# コード解析ツール

コードベースを分析し、複雑度、コメント率、コードクローンなどの重要なメトリクスを収集・可視化するツールです。

## 概要

code-analyzer は、様々なプログラミング言語のソースコードを解析し、以下のようなメトリクスを収集します：

- コード行数の分布（コード行、コメント行、空行）
- サイクロマティック複雑度
- 関数/メソッドのサイズと数
- コメントの割合
- コードクローン（重複コード）の検出
- 時系列メトリクスの変化傾向

これらのメトリクスは、コードの品質評価、リファクタリングの優先順位付け、コードベースの健全性監視に役立ちます。

## インストール

```bash
go install github.com/yourusername/code-analyzer@latest
```

または、リポジトリをクローンして手動でビルドすることもできます：

```bash
git clone https://github.com/yourusername/code-analyzer.git
cd code-analyzer
go build -o code-analyzer ./cmd/code-analyzer
```

## 使用方法

### コマンドラインオプション

以下のオプションを指定できます：

| オプション | 説明 | デフォルト値 |
|------------|------|------------|
| `-path` | 分析対象のプロジェクトパス | `.` |
| `-ext` | 分析対象の拡張子 (.go,.py,.js など) | `.go` |
| `-format` | 出力形式 (text, json, csv) | `text` |
| `-output` | 出力ファイルのパス | - (標準出力) |
| `-visual` | ビジュアルHTMLレポートを生成 | `false` |
| `-history` | 履歴データファイルのパス | - |
| `-detect-clones` | コードクローン検出を有効化 | `false` |
| `-min-block-size` | クローン検出の最小ブロックサイズ | `30` |
| `-min-similarity` | クローン検出の最小類似度 (0.0-1.0) | `0.8` |
| `-v` | 詳細なログを出力 | `false` |

## 使用例

### 単一言語のプロジェクト分析

```bash
code-analyzer -path ./my-project
```

### 複数言語のプロジェクト分析

```bash
code-analyzer -path ./my-project -ext .go,.py,.js
```

### JSONフォーマットでの出力

```bash
code-analyzer -path ./my-project -format json -output metrics.json
```

### ビジュアルレポートの生成

```bash
code-analyzer -path ./my-project -visual -output report
```

### コードクローン検出の有効化

```bash
code-analyzer -path ./my-project -detect-clones
```

### クローン検出のパラメータ調整

```bash
code-analyzer -path ./my-project -detect-clones -min-block-size 20 -min-similarity 0.7
```

### 時系列分析の実行

```bash
code-analyzer -path ./my-project -history metrics_history.json
```

## 出力形式

### テキスト形式

```
Project Analysis for: ./my-project
Analyzed at: 2025-05-10 15:30:45

Files analyzed: 42
Total lines: 8547
  - Code lines: 5637 (65.95%)
  - Comment lines: 1245 (14.57%)
  - Blank lines: 1665 (19.48%)
Average complexity: 4.32
Maximum complexity: 24
Comment-to-code ratio: 22.09%

Files with highest complexity:
  src/controllers/user_controller.go: 24
  src/services/auth_service.go: 17
  src/utils/validator.go: 15
  src/models/product.go: 12
  src/middleware/auth.go: 11
```

### JSON形式

```json
{
  "project_path": "./my-project",
  "analyzed_at": "2025-05-10T15:30:45Z",
  "file_count": 42,
  "total_lines": 8547,
  "total_code_lines": 5637,
  "total_comments": 1245,
  "total_blank_lines": 1665,
  "avg_complexity": 4.32,
  "max_complexity": 24,
  "comment_ratio": 22.09,
  "files": [
    {
      "path": "src/controllers/user_controller.go",
      "total_lines": 320,
      "code_lines": 245,
      "comment_lines": 40,
      "blank_lines": 35,
      "function_count": 8,
      "avg_function_size": 27.5,
      "max_function_size": 65,
      "complexity": 24,
      "comment_ratio": 16.33
    },
    // ...
  ],
  "clones": [
    {
      "source_file": "src/utils/logger.go",
      "target_file": "src/utils/debugger.go",
      "source_line": 45,
      "target_line": 23,
      "line_count": 12,
      "similarity": 0.92,
      "code": "func formatLogMessage(level, message string) string { ... }"
    },
    // ...
  ]
}
```

### CSV形式

```
Path,TotalLines,CodeLines,CommentLines,BlankLines,FunctionCount,AvgFunctionSize,MaxFunctionSize,Complexity,CommentRatio
src/controllers/user_controller.go,320,245,40,35,8,27.5,65,24,16.33
src/services/auth_service.go,285,210,45,30,6,31.2,58,17,21.43
...

SUMMARY,42 files,8547 lines,5637 code,1245 comments,1665 blank,4.32 avg complexity,24 max complexity,22.09% comment ratio
```

### ビジュアルHTMLレポート

ビジュアルレポートは以下のような情報を含む対話型HTMLページを生成します：

- メトリクスのサマリーカード
- コード行分布の円グラフ
- 複雑度分布のヒストグラム
- 時系列メトリクスのトレンドグラフ
- コードクローンの分析
- 複雑度に基づいて色分けされたファイル一覧表

## 特徴

### コード複雑度分析

コードの複雑さをサイクロマティック複雑度で測定します。条件分岐、ループ、論理演算子などに基づいて計算され、テストの難易度やバグの潜在的なリスクを示します。

### コメント率の分析

コード行に対するコメント行の比率を計算します。適切なドキュメンテーションの指標となります。

### コードクローン検出

プロジェクト内の重複コードを検出します。リファクタリングの機会を特定し、コードの保守性を向上させるのに役立ちます。

### 時系列分析

過去の分析結果と比較して、メトリクスの経時変化を追跡します。コードベースの健全性の改善または悪化を監視できます。

### ビジュアルレポート

メトリクスやコードクローンを視覚的に表現したHTMLレポートを生成します。データをより直感的に理解し、チームメンバーと共有しやすくなります。

## サポートされる言語

基本的なメトリクス（行数、コメント率）はすべてのテキストベースの言語に対応しています。

詳細な解析（複雑度、関数解析）は以下の言語で利用可能です：

- Go (`.go`)

他の言語については、基本的なメトリクスのみサポートされます。将来のリリースで詳細分析のサポートが拡張される予定です。

## 活用シナリオ

### コード品質の監視

定期的な分析を実行してコードベースの品質メトリクスを追跡し、問題が深刻化する前に対処します。

```bash
code-analyzer -path ./my-project -history metrics_history.json -visual -output reports/$(date +%Y%m%d)
```

### リファクタリングの優先順位付け

複雑度の高いファイルやコードクローンを特定し、リファクタリングの優先順位を決定します。

```bash
code-analyzer -path ./my-project -detect-clones -format json -output refactoring-targets.json
```

### CIパイプラインでの品質チェック

CI/CDパイプラインで使用して、コードの品質基準を維持します。

```bash
code-analyzer -path ./my-project -format json | jq '.avg_complexity < 5 and .comment_ratio > 15'
```
