package usecases

import (
	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
	checkbodylength "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/check_body_length"
	craftmarkdown "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/craft_markdown"
	distributefiles "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/distribute_files"
	grepstr "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/grep_str"
	migratetomemos "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/migrate_to_memos"
	renamebodiesbycategoryid "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/rename_bodies_by_category_id"
)

type Service struct {
	distributeFilesOperation distributeFilesOperation
	craftMarkdownOperation   craftMarkdownOperation
	checkBodyLengthOperation checkBodyLengthOperation
	grepStrOperation         grepStrOperation
	renameBodiesOperation    renameBodiesByCategoryIDOperation
	migrateToMemosOperation  migrateToMemosOperation
}

func NewService(fileSystem filesystem.Repository) *Service {
	repo := fileSystem
	if repo == nil {
		repo = filesystem.NewRepository()
	}

	return newServiceWithOperations(
		distributefiles.NewService(repo),
		craftmarkdown.NewService(repo),
		checkbodylength.NewService(repo),
		grepstr.NewService(repo),
		renamebodiesbycategoryid.NewService(repo),
		migratetomemos.NewService(repo, nil),
	)
}
