package usecases

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	youtube "github.com/kkdai/youtube/v2"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// YouTubeClientImpl はYouTubeClientの実装です
type YouTubeClientImpl struct {
	client *youtube.Client
}

// NewYouTubeClientImpl は新しいYouTubeClientImplを作成します
func NewYouTubeClientImpl() *YouTubeClientImpl {
	client := &youtube.Client{
		HTTPClient:  http.DefaultClient,
		MaxRoutines: 10,
		ChunkSize:   10 * 1024 * 1024, // 10MB
	}
	return &YouTubeClientImpl{
		client: client,
	}
}

// NewYouTubeClientImplWithClient は指定されたクライアントでYouTubeClientImplを作成します
func NewYouTubeClientImplWithClient(client *youtube.Client) *YouTubeClientImpl {
	return &YouTubeClientImpl{
		client: client,
	}
}

// SetMaxRoutines はMaxRoutinesを設定します
func (y *YouTubeClientImpl) SetMaxRoutines(maxRoutines int) {
	y.client.MaxRoutines = maxRoutines
}

// SetChunkSize はChunkSizeを設定します
func (y *YouTubeClientImpl) SetChunkSize(chunkSize int64) {
	y.client.ChunkSize = chunkSize
}

func (y *YouTubeClientImpl) GetVideoContext(ctx context.Context, url string) (*youtube.Video, error) {
	return y.client.GetVideoContext(ctx, url)
}

func (y *YouTubeClientImpl) GetPlaylistContext(ctx context.Context, url string) (*youtube.Playlist, error) {
	return y.client.GetPlaylistContext(ctx, url)
}

func (y *YouTubeClientImpl) VideoFromPlaylistEntryContext(ctx context.Context, entry *youtube.PlaylistEntry) (*youtube.Video, error) {
	return y.client.VideoFromPlaylistEntryContext(ctx, entry)
}

func (y *YouTubeClientImpl) GetStreamContext(ctx context.Context, video *youtube.Video, format *youtube.Format) (io.ReadCloser, int64, error) {
	return y.client.GetStreamContext(ctx, video, format)
}

// FileSystemImpl はFileSystemの実装です
type FileSystemImpl struct{}

// NewFileSystemImpl は新しいFileSystemImplを作成します
func NewFileSystemImpl() *FileSystemImpl {
	return &FileSystemImpl{}
}

func (f *FileSystemImpl) MkdirAll(path string, perm uint32) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

func (f *FileSystemImpl) Create(name string) (File, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &FileImpl{file: file}, nil
}

func (f *FileSystemImpl) Remove(name string) error {
	return os.Remove(name)
}

func (f *FileSystemImpl) Stat(name string) (FileInfo, error) {
	info, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	return &FileInfoImpl{info: info}, nil
}

// FileImpl はFileの実装です
type FileImpl struct {
	file *os.File
}

func (f *FileImpl) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

func (f *FileImpl) Close() error {
	return f.file.Close()
}

func (f *FileImpl) Stat() (FileInfo, error) {
	info, err := f.file.Stat()
	if err != nil {
		return nil, err
	}
	return &FileInfoImpl{info: info}, nil
}

// FileInfoImpl はFileInfoの実装です
type FileInfoImpl struct {
	info os.FileInfo
}

func (f *FileInfoImpl) Size() int64 {
	return f.info.Size()
}

// VideoProcessorImpl はVideoProcessorの実装です
type VideoProcessorImpl struct{}

// NewVideoProcessorImpl は新しいVideoProcessorImplを作成します
func NewVideoProcessorImpl() *VideoProcessorImpl {
	return &VideoProcessorImpl{}
}

func (v *VideoProcessorImpl) MergeVideoAudio(videoPath, audioPath, outputPath string) error {
	videoInput := ffmpeg.Input(videoPath)
	audioInput := ffmpeg.Input(audioPath)

	return ffmpeg.Output([]*ffmpeg.Stream{videoInput, audioInput}, outputPath, ffmpeg.KwArgs{
		"c:v": "copy",                     // 映像をそのままコピー（再エンコードなし）
		"c:a": "aac",                      // 音声をAACでエンコード
		"map": []string{"0:v:0", "1:a:0"}, // 映像と音声のマッピング
	}).OverWriteOutput().Run()
}

// TimeProviderImpl はTimeProviderの実装です
type TimeProviderImpl struct{}

// NewTimeProviderImpl は新しいTimeProviderImplを作成します
func NewTimeProviderImpl() *TimeProviderImpl {
	return &TimeProviderImpl{}
}

func (t *TimeProviderImpl) Now() time.Time {
	return time.Now()
}

func (t *TimeProviderImpl) LoadLocation(name string) (*time.Location, error) {
	return time.LoadLocation(name)
}
