package common

import "os"

const (
	DefaultCategory           = "uncategorized"
	SupportedPageType         = "content"
	SupportedPageTypeContent  = "content"
	SupportedPageTypeArtifact = "artifact"
	DefaultDirectoryPerm      = os.FileMode(0755)
	RequiredBackupTag         = "91-backup/tool-migration/202602-notion"
)
