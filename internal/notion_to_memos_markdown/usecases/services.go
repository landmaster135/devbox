package usecases

import (
	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
	artifactcraftmarkdown "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/artifact_craft_markdown"
	checkbodylength "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/content_check_body_length"
	contentcraftmarkdown "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/content_craft_markdown"
	distributefiles "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/content_distribute_files"
	grepstr "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/content_grep_str"
	migratetomemos "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/content_migrate_to_memos"
	renamebodiesbycategoryid "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/operations/content_rename_bodies_by_category_id"
)

type Service struct {
	distributeFilesOperation distributeFilesOperation
	craftMarkdownOperation   craftMarkdownOperation
	artifactCraftOperation   craftMarkdownOperation
	checkBodyLengthOperation checkBodyLengthOperation
	grepStrOperation         grepStrOperation
	renameBodiesOperation    renameBodiesByCategoryIDOperation
	migrateToMemosOperation  migrateToMemosOperation
}

func NewService(fileSystem filesystem.Repository) *Service {
	return NewServiceWithReporter(fileSystem, nil)
}

func NewServiceWithReporter(fileSystem filesystem.Repository, reporter migratetomemos.ProgressReporter) *Service {
	repo := fileSystem
	if repo == nil {
		repo = filesystem.NewRepository()
	}

	return newServiceWithOperations(
		distributefiles.NewService(repo),
		contentcraftmarkdown.NewService(repo),
		artifactcraftmarkdown.NewService(repo),
		checkbodylength.NewService(repo),
		grepstr.NewService(repo),
		renamebodiesbycategoryid.NewService(repo),
		migratetomemos.NewService(repo, nil, reporter),
	)
}
