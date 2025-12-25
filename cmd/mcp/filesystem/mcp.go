package filesystem

import (
	"context"
	"fmt"
	"strings"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/filesystem/usecases"
)

// #==============================================================#
// ##          Handlers                                          ##
// #==============================================================#
func handleReadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return nil, err
	}

	// サービスの初期化
	fsService := usecases.NewFileSystemService([]string{path})
	offset := request.GetInt("offset", 1)
	limit := request.GetInt("limit", 2000)
	content, err := fsService.ReadFile(path, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み取りに失敗しました: %v", err)
	}

	return mcp.NewToolResultText(content), nil
}

func handleWriteFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return nil, err
	}

	content, err := request.RequireString("content")
	if err != nil {
		return nil, err
	}

	// サービスの初期化
	fsService := usecases.NewFileSystemService([]string{path})
	err = fsService.WriteFile(path, content)
	if err != nil {
		return nil, fmt.Errorf("ファイルの書き込みに失敗しました: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("ファイル %s への書き込みに成功しました", path)), nil
}

func handleCreateDirectory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return nil, err
	}

	// サービスの初期化
	fsService := usecases.NewFileSystemService([]string{path})
	err = fsService.CreateDirectory(path)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの作成に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("ディレクトリ %s の作成に成功しました", path)), nil
}

func handleListDirectory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return nil, err
	}
	offset := request.GetInt("offset", 1)
	limit := request.GetInt("limit", 25)
	if offset <= 0 {
		return nil, fmt.Errorf("offsetは1以上で指定してください")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("limitは1以上で指定してください")
	}

	// サービスの初期化
	fsService := usecases.NewFileSystemService([]string{path})
	entries, err := fsService.ListDirectoryWithOptions(path, usecases.ListDirectoryOptions{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの一覧取得に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(strings.Join(entries, "\n")), nil
}

func handleDirectoryTree(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return nil, err
	}
	offset := request.GetInt("offset", 1)
	limit := request.GetInt("limit", 25)
	depth := request.GetInt("depth", 0)
	if offset <= 0 {
		return nil, fmt.Errorf("offsetは1以上で指定してください")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("limitは1以上で指定してください")
	}
	if depth < 0 {
		return nil, fmt.Errorf("depthは0以上で指定してください")
	}

	// サービスの初期化
	fsService := usecases.NewFileSystemService([]string{path})
	yamlStr, err := fsService.GetDirectoryTreeAsYAMLWithOptions(path, usecases.DirectoryTreeOptions{
		Offset: offset,
		Limit:  limit,
		Depth:  depth,
	})
	if err != nil {
		return nil, fmt.Errorf("ディレクトリツリーの取得に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(yamlStr), nil
}

func handleMoveFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	source, err := request.RequireString("source")
	if err != nil {
		return nil, err
	}

	destination, err := request.RequireString("destination")
	if err != nil {
		return nil, err
	}

	// サービスの初期化（両方のパスを許可）
	fsService := usecases.NewFileSystemService([]string{source, destination})
	err = fsService.MoveFile(source, destination)
	if err != nil {
		return nil, fmt.Errorf("ファイルの移動に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("%s から %s への移動に成功しました", source, destination)), nil
}

func handleSearchFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return nil, err
	}

	pattern, err := request.RequireString("pattern")
	if err != nil {
		return nil, err
	}

	var excludePatterns []string
	excludePattern := request.GetString("exclude_pattern", "")
	if excludePattern != "" {
		excludePatterns = append(excludePatterns, excludePattern)
	}

	// サービスの初期化
	fsService := usecases.NewFileSystemService([]string{path})
	results, err := fsService.SearchFiles(path, pattern, excludePatterns)
	if err != nil {
		return nil, fmt.Errorf("ファイルの検索に失敗しました: %v", err)
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("一致するものが見つかりませんでした"), nil
	}

	return mcp.NewToolResultText(strings.Join(results, "\n")), nil
}

func handleGetFileInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return nil, err
	}

	// サービスの初期化
	fsService := usecases.NewFileSystemService([]string{path})
	result, err := fsService.GetFileInfoAsText(path)
	if err != nil {
		return nil, fmt.Errorf("ファイル情報の取得に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

func handleListAllowedDirectories(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return nil, err
	}

	// サービスの初期化
	fsService := usecases.NewFileSystemService([]string{path})
	result := fsService.GetAllowedDirectoriesAsText()

	return mcp.NewToolResultText(result), nil
}

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt("system_prompt_01",
		mcp.WithPromptDescription("This is a prompt for the filesystem."),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for the filesystem.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You use this great filesystem well."),
				},
			},
		}, nil
	})
	return s
}

// BuildFileSystemServer はファイルシステムMCPサーバーを構築する関数です
func BuildFileSystemServer() {
	// サーバーの設定
	s := server.NewMCPServer(
		"Secure Filesystem Server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)

	// ファイル読み取りツール
	readFileTool := mcp.NewTool("read_file",
		mcp.WithDescription("ファイルシステムからファイルの内容を読み取ります。様々なテキストエンコーディングを処理し、ファイルが読み取れない場合は詳細なエラーメッセージを提供します。単一のファイルの内容を調べる必要がある場合に使用します。許可されたディレクトリ内でのみ動作します。"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("読み取るファイルのパス"),
		),
		mcp.WithNumber("offset",
			mcp.Description("開始行番号。省略時は1"),
		),
		mcp.WithNumber("limit",
			mcp.Description("始まりから返す最大行数。省略時は2000"),
		),
	)

	s.AddTool(readFileTool, handleReadFile)

	// ファイル書き込みツール
	writeFileTool := mcp.NewTool("write_file",
		mcp.WithDescription("新しいファイルを作成するか、既存のファイルを新しい内容で完全に上書きします。警告なしに既存のファイルを上書きするため、注意して使用してください。適切なエンコーディングでテキスト内容を処理します。許可されたディレクトリ内でのみ動作します。"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("書き込むファイルのパス"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("ファイルに書き込む内容"),
		),
	)

	s.AddTool(writeFileTool, handleWriteFile)

	// ディレクトリ作成ツール
	createDirTool := mcp.NewTool("create_directory",
		mcp.WithDescription("新しいディレクトリを作成するか、ディレクトリが存在することを確認します。1回の操作で複数のネストされたディレクトリを作成できます。ディレクトリが既に存在する場合、この操作は静かに成功します。プロジェクトのディレクトリ構造を設定したり、必要なパスが存在することを確認したりするのに最適です。許可されたディレクトリ内でのみ動作します。"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("作成するディレクトリのパス"),
		),
	)

	s.AddTool(createDirTool, handleCreateDirectory)

	// ディレクトリ一覧ツール
	listDirTool := mcp.NewTool("list_directory",
		mcp.WithDescription("指定されたパス内のすべてのファイルとディレクトリの詳細な一覧を取得します。結果は[FILE]と[DIR]のプレフィックスでファイルとディレクトリを明確に区別します。このツールはディレクトリ構造を理解し、ディレクトリ内の特定のファイルを見つけるのに不可欠です。許可されたディレクトリ内でのみ動作します。"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("一覧表示するディレクトリのパス"),
		),
		mcp.WithNumber("offset",
			mcp.Description("結果の開始ディレクトリの位置。省略時は1"),
		),
		mcp.WithNumber("limit",
			mcp.Description("返す最大件数。省略時は25"),
		),
	)

	s.AddTool(listDirTool, handleListDirectory)

	// ディレクトリツリーツール
	dirTreeTool := mcp.NewTool("directory_tree",
		mcp.WithDescription("ファイルとディレクトリの再帰的なツリービューをJSON構造として取得します。各エントリには「name」、「type」（file/directory）、ディレクトリの場合は「children」が含まれます。ファイルには子配列がなく、ディレクトリには常に子配列（空の場合もあります）があります。出力は読みやすさのために2スペースのインデントでフォーマットされています。許可されたディレクトリ内でのみ動作します。"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("ツリーを取得するディレクトリのパス"),
		),
		mcp.WithNumber("offset",
			mcp.Description("開始するディレクトリの位置。省略時は1"),
		),
		mcp.WithNumber("limit",
			mcp.Description("返す最大件数。省略時は25"),
		),
		mcp.WithNumber("depth",
			mcp.Description("辿る最大深さ。0は無制限"),
		),
	)

	s.AddTool(dirTreeTool, handleDirectoryTree)

	// ファイル移動ツール
	moveFileTool := mcp.NewTool("move_file",
		mcp.WithDescription("ファイルとディレクトリを移動または名前変更します。ディレクトリ間でファイルを移動し、1回の操作でそれらの名前を変更できます。宛先が存在する場合、操作は失敗します。異なるディレクトリ間で動作し、同じディレクトリ内での単純な名前変更にも使用できます。ソースと宛先の両方が許可されたディレクトリ内にある必要があります。"),
		mcp.WithString("source",
			mcp.Required(),
			mcp.Description("移動するファイルまたはディレクトリのパス"),
		),
		mcp.WithString("destination",
			mcp.Required(),
			mcp.Description("移動先のパス"),
		),
	)

	s.AddTool(moveFileTool, handleMoveFile)

	// ファイル検索ツール
	searchFilesTool := mcp.NewTool("search_files",
		mcp.WithDescription("パターンに一致するファイルとディレクトリを再帰的に検索します。開始パスからすべてのサブディレクトリを検索します。検索は大文字と小文字を区別せず、部分的な名前に一致します。一致するすべての項目への完全なパスを返します。正確な場所がわからないファイルを見つけるのに最適です。許可されたディレクトリ内でのみ検索します。"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("検索を開始するディレクトリのパス"),
		),
		mcp.WithString("pattern",
			mcp.Required(),
			mcp.Description("検索パターン"),
		),
		mcp.WithString("exclude_pattern",
			mcp.Description("除外するパターン（オプション）"),
		),
	)

	s.AddTool(searchFilesTool, handleSearchFiles)

	// ファイル情報取得ツール
	fileInfoTool := mcp.NewTool("get_file_info",
		mcp.WithDescription("ファイルまたはディレクトリに関する詳細なメタデータを取得します。サイズ、作成時間、最終変更時間、権限、タイプなどの包括的な情報を返します。このツールは、実際の内容を読み取らずにファイルの特性を理解するのに最適です。許可されたディレクトリ内でのみ動作します。"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("情報を取得するファイルまたはディレクトリのパス"),
		),
	)

	s.AddTool(fileInfoTool, handleGetFileInfo)

	// 許可されたディレクトリ一覧ツール
	allowedDirsTool := mcp.NewTool("list_allowed_directories",
		mcp.WithDescription("このサーバーがアクセスできるディレクトリのリストを返します。ファイルにアクセスする前に、どのディレクトリが利用可能かを理解するために使用します。"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("確認するパス"),
		),
	)

	s.AddTool(allowedDirsTool, handleListAllowedDirectories)

	// プロンプト
	s = addPromptIntoServer(s)

	// サーバーの起動
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("サーバーエラー: %v\n", err)
	}
}
