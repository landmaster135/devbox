package memos

import "context"

type CreateMemoCall struct {
	Content string
}

type PatchFilesCall struct {
	Memo      string
	FilePaths []string
}

type MockClient struct {
	CreateMemoFunc func(ctx context.Context, content string) (string, error)
	PatchFilesFunc func(ctx context.Context, memo string, filePaths []string) error

	CreateMemoCalls []CreateMemoCall
	PatchFilesCalls []PatchFilesCall
}

func (m *MockClient) CreateMemo(ctx context.Context, content string) (string, error) {
	m.CreateMemoCalls = append(m.CreateMemoCalls, CreateMemoCall{
		Content: content,
	})
	if m.CreateMemoFunc != nil {
		return m.CreateMemoFunc(ctx, content)
	}
	return "", nil
}

func (m *MockClient) PatchFiles(ctx context.Context, memo string, filePaths []string) error {
	clonedPaths := append([]string(nil), filePaths...)
	m.PatchFilesCalls = append(m.PatchFilesCalls, PatchFilesCall{
		Memo:      memo,
		FilePaths: clonedPaths,
	})
	if m.PatchFilesFunc != nil {
		return m.PatchFilesFunc(ctx, memo, clonedPaths)
	}
	return nil
}
