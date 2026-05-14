package config

import (
	"os"
	"strings"
	"testing"
)

func TestConfig_ParseFlags_CreateMemo_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=create-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content=hello",
		"-visibility=private",
		"-state=normal",
		"-pinned",
		"-display-time=2026-02-12T10:00:00Z",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}

	if cfg.Operation != OperationCreateMemo {
		t.Fatalf("operation = %s, want %s", cfg.Operation, OperationCreateMemo)
	}
	if cfg.Visibility != "PRIVATE" {
		t.Fatalf("visibility = %s, want PRIVATE", cfg.Visibility)
	}
	if cfg.State != "NORMAL" {
		t.Fatalf("state = %s, want NORMAL", cfg.State)
	}
	if !cfg.PinnedSet || !cfg.Pinned {
		t.Fatalf("pinned = %v (set=%v), want true and set", cfg.Pinned, cfg.PinnedSet)
	}
}

func TestConfig_ParseFlags_CreateMemoWithContentFile_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=create-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content-file=/tmp/memo.md",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if cfg.ContentFile != "/tmp/memo.md" {
		t.Fatalf("content-file = %s, want /tmp/memo.md", cfg.ContentFile)
	}
}

func TestConfig_ParseFlags_CreateMemoWithContentAndFile_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=create-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content=hello",
		"-content-file=/tmp/memo.md",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "同時に指定できません") {
		t.Fatalf("error = %v, want 同時に指定できません", err)
	}
}

func TestConfig_ParseFlags_UpdateMemoRequiresParams(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=update-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo パラメータ") {
		t.Fatalf("error = %v, want memo パラメータ", err)
	}
}

func TestConfig_ParseFlags_ListMemosPageSizeZero_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=0",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "page-size") {
		t.Fatalf("error = %v, want page-size", err)
	}
}

func TestConfig_ParseFlags_ListMemosWithFilter_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		`-filter=  created_ts > "2023-01-01T13:00:00Z" && visibility == "PUBLIC"  `,
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}

	want := `created_ts > "2023-01-01T13:00:00Z" && visibility == "PUBLIC"`
	if cfg.Filter != want {
		t.Fatalf("filter = %q, want %q", cfg.Filter, want)
	}
}

func TestConfig_ParseFlags_ListMemosWithAnyTags_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-any-tags= health, book ,,",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}

	if cfg.AnyTags != "health, book ,," {
		t.Fatalf("anyTags = %q, want %q", cfg.AnyTags, "health, book ,,")
	}
}

func TestConfig_ParseFlags_ListMemosWithAnyTagsEmptyCSV_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-any-tags= ,, ",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "any-tags に少なくとも1つのタグ") {
		t.Fatalf("error = %v, want any-tags validation error", err)
	}
}

func TestConfig_ParseFlags_ListMemosWithAllTags_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-all-tags= health, book ,,",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}

	if cfg.AllTags != "health, book ,," {
		t.Fatalf("allTags = %q, want %q", cfg.AllTags, "health, book ,,")
	}
}

func TestConfig_ParseFlags_ListMemosWithAllTagsEmptyCSV_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-all-tags= ,, ",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "all-tags に少なくとも1つのタグ") {
		t.Fatalf("error = %v, want all-tags validation error", err)
	}
}

func TestConfig_ParseFlags_ListMemosWithAnyAndAllTags_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-any-tags=health",
		"-all-tags=book",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "any-tags と all-tags は同時に指定できません") {
		t.Fatalf("error = %v, want simultaneous tag options error", err)
	}
}

func TestConfig_ParseFlags_ListMemosWithExcludedTags_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-excluded-tags= health, book ,,",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}

	if cfg.ExcludedTags != "health, book ,," {
		t.Fatalf("excludedTags = %q, want %q", cfg.ExcludedTags, "health, book ,,")
	}
}

func TestConfig_ParseFlags_ListMemosWithExcludedTagsEmptyCSV_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-excluded-tags= ,, ",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "excluded-tags に少なくとも1つのタグ") {
		t.Fatalf("error = %v, want excluded-tags validation error", err)
	}
}

func TestConfig_ParseFlags_ListMemosWithAnyContents_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-any-contents= meeting, study ,,",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}

	if cfg.AnyContents != "meeting, study ,," {
		t.Fatalf("anyContents = %q, want %q", cfg.AnyContents, "meeting, study ,,")
	}
}

func TestConfig_ParseFlags_ListMemosWithAllContents_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-all-contents= meeting, study ,,",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}

	if cfg.AllContents != "meeting, study ,," {
		t.Fatalf("allContents = %q, want %q", cfg.AllContents, "meeting, study ,,")
	}
}

func TestConfig_ParseFlags_ListMemosWithAnyAndAllContents_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-any-contents=meeting",
		"-all-contents=study",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "any-contents と all-contents は同時に指定できません") {
		t.Fatalf("error = %v, want simultaneous options error", err)
	}
}

func TestConfig_ParseFlags_ListMemosWithAllContentsEmptyCSV_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-all-contents=  ,,  ",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "all-contents に少なくとも1つのキーワード") {
		t.Fatalf("error = %v, want all-contents validation error", err)
	}
}

func TestConfig_ParseFlags_ListAttachments_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=list-attachments",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=25",
		"-page-token=next-1",
		"-order-by=create_time desc",
		`-filter=memo == "memos/memo-1"`,
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}

	if cfg.Operation != OperationListAttachments {
		t.Fatalf("operation = %s, want %s", cfg.Operation, OperationListAttachments)
	}
	if cfg.PageSize != 25 {
		t.Fatalf("pageSize = %d, want 25", cfg.PageSize)
	}
	if cfg.PageToken != "next-1" {
		t.Fatalf("pageToken = %s, want next-1", cfg.PageToken)
	}
	if cfg.OrderBy != "create_time desc" {
		t.Fatalf("orderBy = %s, want create_time desc", cfg.OrderBy)
	}
	if cfg.Filter != `memo == "memos/memo-1"` {
		t.Fatalf("filter = %q, want memo filter", cfg.Filter)
	}
}

func TestConfig_ParseFlags_ListAttachmentsPageSizeZero_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=list-attachments",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=0",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "page-size") {
		t.Fatalf("error = %v, want page-size", err)
	}
}

func TestConfig_ParseFlags_InvalidOperation_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=unknown",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "未対応") {
		t.Fatalf("error = %v, want 未対応", err)
	}
}

func TestConfig_ParseFlags_InvalidVisibility_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=create-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content=hello",
		"-visibility=team-only",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "visibility") {
		t.Fatalf("error = %v, want visibility", err)
	}
}

func TestConfig_ParseFlags_Help_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{"-help"})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if !cfg.Help {
		t.Fatalf("help = %v, want true", cfg.Help)
	}
}

func TestConfig_ParseFlags_Normal(t *testing.T) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	os.Args = []string{
		"memos-cli",
		"-operation=get-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
	}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}
	if cfg.Operation != OperationGetMemo {
		t.Fatalf("operation = %s, want %s", cfg.Operation, OperationGetMemo)
	}
}

func TestConfig_ParseFlagsFromArgs_MissingBaseURL_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=get-memo",
		"-api-token=test-token",
		"-memo=memo-1",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "base-url") {
		t.Fatalf("error = %v, want base-url", err)
	}
}

func TestConfig_ParseFlagsFromArgs_MissingAPIToken_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=get-memo",
		"-base-url=https://memos.example.com",
		"-memo=memo-1",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "api-token") {
		t.Fatalf("error = %v, want api-token", err)
	}
}

func TestConfig_ParseFlagsFromArgs_InvalidState_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-state=done",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "state") {
		t.Fatalf("error = %v, want state", err)
	}
}

func TestConfig_ParseFlagsFromArgs_InvalidTimeout_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-timeout=0",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("error = %v, want timeout", err)
	}
}

func TestConfig_ParseFlagsFromArgs_NegativePageSize_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=-1",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "page-size") {
		t.Fatalf("error = %v, want page-size", err)
	}
}

func TestConfig_ParseFlagsFromArgs_GetMemoMissingMemo_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=get-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo パラメータ") {
		t.Fatalf("error = %v, want memo パラメータ", err)
	}
}

func TestConfig_ParseFlagsFromArgs_DeleteMemo_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=delete-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-force=true",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if cfg.Operation != OperationDeleteMemo {
		t.Fatalf("operation = %s, want %s", cfg.Operation, OperationDeleteMemo)
	}
	if cfg.Memo != "memo-1" {
		t.Fatalf("memo = %s, want memo-1", cfg.Memo)
	}
	if !cfg.Force {
		t.Fatalf("force = %v, want true", cfg.Force)
	}
}

func TestConfig_ParseFlagsFromArgs_DeleteMemoMissingMemo_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=delete-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo パラメータ") {
		t.Fatalf("error = %v, want memo パラメータ", err)
	}
}

func TestConfig_ParseFlagsFromArgs_DeleteMemoDefaultForceFalse_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=delete-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if cfg.Force {
		t.Fatalf("force = %v, want false", cfg.Force)
	}
}

func TestConfig_ParseFlagsFromArgs_UpdateMemoMissingContent_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=update-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "content または content-file") {
		t.Fatalf("error = %v, want content または content-file", err)
	}
}

func TestConfig_ParseFlagsFromArgs_UpdateMemoWithContentFile_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=update-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-content-file=/tmp/memo.md",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if cfg.ContentFile != "/tmp/memo.md" {
		t.Fatalf("content-file = %s, want /tmp/memo.md", cfg.ContentFile)
	}
}

func TestConfig_ParseFlagsFromArgs_UpdateMemoUpdatesTimeDefaultFalse_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=update-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-content=updated",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if cfg.UpdatesTime {
		t.Fatalf("updatesTime = %v, want false", cfg.UpdatesTime)
	}
}

func TestConfig_ParseFlagsFromArgs_UpdateMemoUpdatesTimeTrue_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=update-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-content=updated",
		"-updates-time=true",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if !cfg.UpdatesTime {
		t.Fatalf("updatesTime = %v, want true", cfg.UpdatesTime)
	}
}

func TestConfig_ParseFlagsFromArgs_UpdateTag_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=update-tag",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-src-tag=#work",
		"-dest-tag=#project",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if cfg.Operation != OperationUpdateTag {
		t.Fatalf("operation = %s, want %s", cfg.Operation, OperationUpdateTag)
	}
	if cfg.SrcTag != "work" {
		t.Fatalf("srcTag = %s, want work", cfg.SrcTag)
	}
	if cfg.DestTag != "project" {
		t.Fatalf("destTag = %s, want project", cfg.DestTag)
	}
}

func TestConfig_ParseFlagsFromArgs_UpdateTagMissingSrcTag_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=update-tag",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-dest-tag=project",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "src-tag") {
		t.Fatalf("error = %v, want src-tag", err)
	}
}

func TestConfig_ParseFlagsFromArgs_UpdateTagMissingDestTag_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=update-tag",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-src-tag=work",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "dest-tag") {
		t.Fatalf("error = %v, want dest-tag", err)
	}
}

func TestConfig_ParseFlagsFromArgs_PatchFiles_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=patch-files",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-files=./a.png,./b.txt",
		"-replaces=true",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if cfg.Operation != OperationPatchFiles {
		t.Fatalf("operation = %s, want %s", cfg.Operation, OperationPatchFiles)
	}
	if cfg.Files != "./a.png,./b.txt" {
		t.Fatalf("files = %s, want ./a.png,./b.txt", cfg.Files)
	}
	if !cfg.Replaces {
		t.Fatalf("replaces = %v, want true", cfg.Replaces)
	}
}

func TestConfig_ParseFlagsFromArgs_PatchFilesDefaultReplaces_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=patch-files",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-files=./a.png",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if cfg.Replaces {
		t.Fatalf("replaces = %v, want false", cfg.Replaces)
	}
}

func TestConfig_ParseFlagsFromArgs_PatchFilesMissingMemo_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=patch-files",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-files=./a.png",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo パラメータ") {
		t.Fatalf("error = %v, want memo パラメータ", err)
	}
}

func TestConfig_ParseFlagsFromArgs_PatchFilesMissingFiles_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=patch-files",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-files=, ,",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "files パラメータ") {
		t.Fatalf("error = %v, want files パラメータ", err)
	}
}

func TestConfig_ParseFlagsFromArgs_ListMemoRelations_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=list-memo-relations",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if cfg.Operation != OperationListMemoRelations {
		t.Fatalf("operation = %s, want %s", cfg.Operation, OperationListMemoRelations)
	}
	if cfg.Memo != "memo-1" {
		t.Fatalf("memo = %s, want memo-1", cfg.Memo)
	}
}

func TestConfig_ParseFlagsFromArgs_ListMemoRelationsMissingMemo_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=list-memo-relations",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo パラメータ") {
		t.Fatalf("error = %v, want memo パラメータ", err)
	}
}

func TestConfig_ParseFlagsFromArgs_AddMemoRelations_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=add-memo-relations",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-related-memos=memo-2,memo-3",
		"-replaces=true",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if cfg.Operation != OperationAddMemoRelations {
		t.Fatalf("operation = %s, want %s", cfg.Operation, OperationAddMemoRelations)
	}
	if cfg.RelatedMemos != "memo-2,memo-3" {
		t.Fatalf("relatedMemos = %s, want memo-2,memo-3", cfg.RelatedMemos)
	}
	if !cfg.Replaces {
		t.Fatalf("replaces = %v, want true", cfg.Replaces)
	}
}

func TestConfig_ParseFlagsFromArgs_AddMemoRelationsDefaultReplaces_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=add-memo-relations",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-related-memos=memo-2",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if cfg.Replaces {
		t.Fatalf("replaces = %v, want false", cfg.Replaces)
	}
}

func TestConfig_ParseFlagsFromArgs_AddMemoRelationsMissingMemo_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=add-memo-relations",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-related-memos=memo-2",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo パラメータ") {
		t.Fatalf("error = %v, want memo パラメータ", err)
	}
}

func TestConfig_ParseFlagsFromArgs_AddMemoRelationsMissingRelatedMemos_Error(t *testing.T) {
	_, err := ParseFlagsFromArgs([]string{
		"-operation=add-memo-relations",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-related-memos=, ,",
	})
	if err == nil {
		t.Fatal("ParseFlagsFromArgs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "related-memos パラメータ") {
		t.Fatalf("error = %v, want related-memos パラメータ", err)
	}
}

func TestConfig_PrintUsage_Normal(t *testing.T) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()
	os.Args = []string{"memos-cli"}

	PrintUsage()
}

func TestParseBool_Normal(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{in: "true", want: true},
		{in: "false", want: false},
		{in: "YES", want: true},
		{in: "n", want: false},
	}

	for _, tc := range tests {
		got, err := parseBool(tc.in)
		if err != nil {
			t.Fatalf("parseBool(%q) error = %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("parseBool(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestParseBool_Invalid_Error(t *testing.T) {
	_, err := parseBool("invalid")
	if err == nil {
		t.Fatal("parseBool() error = nil, want error")
	}
}
