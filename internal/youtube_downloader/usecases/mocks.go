package usecases

import (
	"context"
	"io"
	"time"

	youtube "github.com/kkdai/youtube/v2"
)

// MockYouTubeClient はYouTubeClientのモック実装です
type MockYouTubeClient struct {
	GetVideoContextFunc                   func(ctx context.Context, url string) (*youtube.Video, error)
	GetPlaylistContextFunc                func(ctx context.Context, url string) (*youtube.Playlist, error)
	VideoFromPlaylistEntryContextFunc     func(ctx context.Context, entry *youtube.PlaylistEntry) (*youtube.Video, error)
	GetStreamContextFunc                  func(ctx context.Context, video *youtube.Video, format *youtube.Format) (io.ReadCloser, int64, error)
	SetMaxRoutinesFunc                    func(maxRoutines int)
	SetChunkSizeFunc                      func(chunkSize int64)
}

func (m *MockYouTubeClient) GetVideoContext(ctx context.Context, url string) (*youtube.Video, error) {
	if m.GetVideoContextFunc != nil {
		return m.GetVideoContextFunc(ctx, url)
	}
	return nil, nil
}

func (m *MockYouTubeClient) GetPlaylistContext(ctx context.Context, url string) (*youtube.Playlist, error) {
	if m.GetPlaylistContextFunc != nil {
		return m.GetPlaylistContextFunc(ctx, url)
	}
	return nil, nil
}

func (m *MockYouTubeClient) VideoFromPlaylistEntryContext(ctx context.Context, entry *youtube.PlaylistEntry) (*youtube.Video, error) {
	if m.VideoFromPlaylistEntryContextFunc != nil {
		return m.VideoFromPlaylistEntryContextFunc(ctx, entry)
	}
	return nil, nil
}

func (m *MockYouTubeClient) GetStreamContext(ctx context.Context, video *youtube.Video, format *youtube.Format) (io.ReadCloser, int64, error) {
	if m.GetStreamContextFunc != nil {
		return m.GetStreamContextFunc(ctx, video, format)
	}
	return nil, 0, nil
}

func (m *MockYouTubeClient) SetMaxRoutines(maxRoutines int) {
	if m.SetMaxRoutinesFunc != nil {
		m.SetMaxRoutinesFunc(maxRoutines)
	}
}

func (m *MockYouTubeClient) SetChunkSize(chunkSize int64) {
	if m.SetChunkSizeFunc != nil {
		m.SetChunkSizeFunc(chunkSize)
	}
}

// MockFileSystem はFileSystemのモック実装です
type MockFileSystem struct {
	MkdirAllFunc func(path string, perm uint32) error
	CreateFunc   func(name string) (File, error)
	RemoveFunc   func(name string) error
	StatFunc     func(name string) (FileInfo, error)
}

func (m *MockFileSystem) MkdirAll(path string, perm uint32) error {
	if m.MkdirAllFunc != nil {
		return m.MkdirAllFunc(path, perm)
	}
	return nil
}

func (m *MockFileSystem) Create(name string) (File, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(name)
	}
	return &MockFile{}, nil
}

func (m *MockFileSystem) Remove(name string) error {
	if m.RemoveFunc != nil {
		return m.RemoveFunc(name)
	}
	return nil
}

func (m *MockFileSystem) Stat(name string) (FileInfo, error) {
	if m.StatFunc != nil {
		return m.StatFunc(name)
	}
	return &MockFileInfo{}, nil
}

// MockFile はFileのモック実装です
type MockFile struct {
	WriteFunc func(p []byte) (n int, err error)
	CloseFunc func() error
	StatFunc  func() (FileInfo, error)
}

func (m *MockFile) Write(p []byte) (n int, err error) {
	if m.WriteFunc != nil {
		return m.WriteFunc(p)
	}
	return len(p), nil
}

func (m *MockFile) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *MockFile) Stat() (FileInfo, error) {
	if m.StatFunc != nil {
		return m.StatFunc()
	}
	return &MockFileInfo{}, nil
}

// MockFileInfo はFileInfoのモック実装です
type MockFileInfo struct {
	SizeFunc func() int64
}

func (m *MockFileInfo) Size() int64 {
	if m.SizeFunc != nil {
		return m.SizeFunc()
	}
	return 0
}

// MockVideoProcessor はVideoProcessorのモック実装です
type MockVideoProcessor struct {
	MergeVideoAudioFunc func(videoPath, audioPath, outputPath string) error
}

func (m *MockVideoProcessor) MergeVideoAudio(videoPath, audioPath, outputPath string) error {
	if m.MergeVideoAudioFunc != nil {
		return m.MergeVideoAudioFunc(videoPath, audioPath, outputPath)
	}
	return nil
}

// MockTimeProvider はTimeProviderのモック実装です
type MockTimeProvider struct {
	NowFunc          func() time.Time
	LoadLocationFunc func(name string) (*time.Location, error)
}

func (m *MockTimeProvider) Now() time.Time {
	if m.NowFunc != nil {
		return m.NowFunc()
	}
	return time.Now()
}

func (m *MockTimeProvider) LoadLocation(name string) (*time.Location, error) {
	if m.LoadLocationFunc != nil {
		return m.LoadLocationFunc(name)
	}
	return time.LoadLocation(name)
}
