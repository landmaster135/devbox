`openai/codex`のリポジトリ内の`codex-rs`ディレクトリ内で実装されているツールを参考に作るためのメモである。
**`read_file`, `list_dir`, `grep_files`ツールの機能は`filesystem`の方に既に実装済み。**

場所と構成

- ファイルシステム系のツールは実装（handler）が core/src/tools/handlers に、LLM へ公開する JSON Schema は core/src/tools/spec.rs にまとまっています。例えば read_file/list_dir/grep_files/apply_patch の登録や JSON 定義は core/src/tools/spec.rs:481-672、実体のロジックは各 handlers 配下にあります。さらに、該当ツールを有効化するかどうかは experimental_supported_tools フラグ経由で制御され（core/src/tools/spec.rs:1059-1084）、test-gpt-5 など一部モデルでのみオンになるよう core/src/models_manager/model_family.rs:239-254 に記載されています。

代表的なファイルシステム系ツール

- read_file — core/src/tools/handlers/read_file.rs:23 にある handler が、絶対パス・1始まりのオフセット・最大行数などを検証したうえでファイル内容を読み出します。slice モードは単純な行範囲を、indentation モードはアンカー行を中心にインデント階層を辿って関連ブロック（兄弟やドキュコメントをオプションで含める）を返し、各行は L{行番号}: 内容 形式で 500 文字にトリムされます。対応する JSON 定義は core/src/tools/spec.rs:531-625 で、file_path・offset・limit に加えて mode と indentation オプションを受け付ける仕様です。
- list_dir — core/src/tools/handlers/list_dir.rs:36 では BFS でサブディレクトリを辿り、1始まりの offset/limit と最大探索 depth を守りながらツリー表示用の Absolute path: ... 行＋インデント付きの各エントリを生成します。名前は 500 文字で切り詰められ、ディレクトリは /、シンボリックリンクは @、未知タイプは ? が末尾に付与されるため種類が一目で分かります。JSON 定義は core/src/tools/spec.rs:629-672 にあり、dir_path（絶対パス）必須で offset/limit/depth を指定できます。
- grep_files — core/src/tools/handlers/grep_files.rs:1 が ripgrep (rg --files-with-matches --sortr=modified) を 30 秒タイムアウト付きで実行するファイル検索ツールです。正規表現パターン必須、--glob に相当する include、検索対象の path、最大ヒット数 limit（1〜2000）を扱い、結果は更新日時順のファイルパス一覧を返します。Schema は core/src/tools/spec.rs:481-528 で定義されています。
- apply_patch — core/src/tools/handlers/apply_patch.rs:1 が LLM から受け取った apply_patch 形式の差分（JSON 経由 or Freeform）を検証し、ファイル更新を実行する唯一のミューテーティブな標準ツールです。codex_apply_patch::maybe_parse_apply_patch_verified で妥当性と対象パスを確認し、直接適用できる場合はapply_patch::apply_patch を呼び出し、ユーザー承認が必要なケースは専用ランタイム（ApplyPatchRuntime）へ委譲します。対応する ToolSpec は core/src/tools/spec.rs 冒頭付近（create_apply_patch_freeform_tool/create_apply_patch_json_tool）にまとまっており、差分本文と適用メタ情報を受け取れるようになっています。
