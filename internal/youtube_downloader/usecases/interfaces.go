package usecases

import (
	"context"
	"io"
	"time"

	youtube "github.com/kkdai/youtube/v2"
)

// YouTubeClient はYouTube APIクライアントのインターフェースです
type YouTubeClient interface {
	GetVideoContext(ctx context.Context, url string) (*youtube.Video, error)
	GetPlaylistContext(ctx context.Context, url string) (*youtube.Playlist, error)
	VideoFromPlaylistEntryContext(ctx context.Context, entry *youtube.PlaylistEntry) (*youtube.Video, error)
	GetStreamContext(ctx context.Context, video *youtube.Video, format *youtube.Format) (io.ReadCloser, int64, error)
	SetMaxRoutines(maxRoutines int)
	SetChunkSize(chunkSize int64)
}

// FileSystem はファイルシステム操作のインターフェースです
type FileSystem interface {
	MkdirAll(path string, perm uint32) error
	Create(name string) (File, error)
	Remove(name string) error
	Stat(name string) (FileInfo, error)
}

// File はファイル操作のインターフェースです
type File interface {
	io.Writer
	io.Closer
	Stat() (FileInfo, error)
}

// FileInfo はファイル情報のインターフェースです
type FileInfo interface {
	Size() int64
}

// VideoProcessor は動画処理のインターフェースです
type VideoProcessor interface {
	MergeVideoAudio(videoPath, audioPath, outputPath string) error
}

// TimeProvider は時刻取得のインターフェースです
type TimeProvider interface {
	Now() time.Time
	LoadLocation(name string) (*time.Location, error)
}
