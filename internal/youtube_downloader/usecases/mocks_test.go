package usecases

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	youtube "github.com/kkdai/youtube/v2"
)

// #==============================================================#
// ##         MockYouTubeClient Tests                           ##
// #==============================================================#
type MockYouTubeClientTestSuite struct {
	mockClient *MockYouTubeClient
}

func (suite *MockYouTubeClientTestSuite) SetupTest() {
	suite.mockClient = &MockYouTubeClient{}
}

// TestMockYouTubeClient_GetVideoContext_Normal は GetVideoContext メソッドの正常系テスト
func TestMockYouTubeClient_GetVideoContext_Normal(t *testing.T) {
	// Arrange
	suite := &MockYouTubeClientTestSuite{}
	suite.SetupTest()
	ctx := context.Background()
	url := "https://www.youtube.com/watch?v=test"
	expectedVideo := &youtube.Video{
		ID:    "test",
		Title: "Test Video",
	}

	suite.mockClient.GetVideoContextFunc = func(ctx context.Context, url string) (*youtube.Video, error) {
		return expectedVideo, nil
	}

	// Act
	result, err := suite.mockClient.GetVideoContext(ctx, url)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedVideo, result)
}

// TestMockYouTubeClient_GetVideoContext_Error は GetVideoContext メソッドのエラー系テスト
func TestMockYouTubeClient_GetVideoContext_Error(t *testing.T) {
	// Arrange
	suite := &MockYouTubeClientTestSuite{}
	suite.SetupTest()
	ctx := context.Background()
	url := "https://www.youtube.com/watch?v=invalid"
	expectedError := errors.New("video not found")

	suite.mockClient.GetVideoContextFunc = func(ctx context.Context, url string) (*youtube.Video, error) {
		return nil, expectedError
	}

	// Act
	result, err := suite.mockClient.GetVideoContext(ctx, url)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

// TestMockYouTubeClient_GetPlaylistContext_Normal は GetPlaylistContext メソッドの正常系テスト
func TestMockYouTubeClient_GetPlaylistContext_Normal(t *testing.T) {
	// Arrange
	suite := &MockYouTubeClientTestSuite{}
	suite.SetupTest()
	ctx := context.Background()
	url := "https://www.youtube.com/playlist?list=test"
	expectedPlaylist := &youtube.Playlist{
		ID:    "test",
		Title: "Test Playlist",
	}

	suite.mockClient.GetPlaylistContextFunc = func(ctx context.Context, url string) (*youtube.Playlist, error) {
		return expectedPlaylist, nil
	}

	// Act
	result, err := suite.mockClient.GetPlaylistContext(ctx, url)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPlaylist, result)
}

// TestMockYouTubeClient_GetPlaylistContext_Error は GetPlaylistContext メソッドのエラー系テスト
func TestMockYouTubeClient_GetPlaylistContext_Error(t *testing.T) {
	// Arrange
	suite := &MockYouTubeClientTestSuite{}
	suite.SetupTest()
	ctx := context.Background()
	url := "https://www.youtube.com/playlist?list=invalid"
	expectedError := errors.New("playlist not found")

	suite.mockClient.GetPlaylistContextFunc = func(ctx context.Context, url string) (*youtube.Playlist, error) {
		return nil, expectedError
	}

	// Act
	result, err := suite.mockClient.GetPlaylistContext(ctx, url)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

// TestMockYouTubeClient_VideoFromPlaylistEntryContext_Normal は VideoFromPlaylistEntryContext メソッドの正常系テスト
func TestMockYouTubeClient_VideoFromPlaylistEntryContext_Normal(t *testing.T) {
	// Arrange
	suite := &MockYouTubeClientTestSuite{}
	suite.SetupTest()
	ctx := context.Background()
	entry := &youtube.PlaylistEntry{
		ID: "test",
	}
	expectedVideo := &youtube.Video{
		ID:    "test",
		Title: "Test Video from Playlist",
	}

	suite.mockClient.VideoFromPlaylistEntryContextFunc = func(ctx context.Context, entry *youtube.PlaylistEntry) (*youtube.Video, error) {
		return expectedVideo, nil
	}

	// Act
	result, err := suite.mockClient.VideoFromPlaylistEntryContext(ctx, entry)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedVideo, result)
}

// TestMockYouTubeClient_VideoFromPlaylistEntryContext_Error は VideoFromPlaylistEntryContext メソッドのエラー系テスト
func TestMockYouTubeClient_VideoFromPlaylistEntryContext_Error(t *testing.T) {
	// Arrange
	suite := &MockYouTubeClientTestSuite{}
	suite.SetupTest()
	ctx := context.Background()
	entry := &youtube.PlaylistEntry{
		ID: "invalid",
	}
	expectedError := errors.New("video from playlist entry not found")

	suite.mockClient.VideoFromPlaylistEntryContextFunc = func(ctx context.Context, entry *youtube.PlaylistEntry) (*youtube.Video, error) {
		return nil, expectedError
	}

	// Act
	result, err := suite.mockClient.VideoFromPlaylistEntryContext(ctx, entry)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

// TestMockYouTubeClient_GetStreamContext_Normal は GetStreamContext メソッドの正常系テスト
func TestMockYouTubeClient_GetStreamContext_Normal(t *testing.T) {
	// Arrange
	suite := &MockYouTubeClientTestSuite{}
	suite.SetupTest()
	ctx := context.Background()
	video := &youtube.Video{ID: "test"}
	format := &youtube.Format{ItagNo: 18}
	expectedStream := io.NopCloser(strings.NewReader("test stream data"))
	expectedSize := int64(1024)

	suite.mockClient.GetStreamContextFunc = func(ctx context.Context, video *youtube.Video, format *youtube.Format) (io.ReadCloser, int64, error) {
		return expectedStream, expectedSize, nil
	}

	// Act
	stream, size, err := suite.mockClient.GetStreamContext(ctx, video, format)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedStream, stream)
	assert.Equal(t, expectedSize, size)
}

// TestMockYouTubeClient_GetStreamContext_Error は GetStreamContext メソッドのエラー系テスト
func TestMockYouTubeClient_GetStreamContext_Error(t *testing.T) {
	// Arrange
	suite := &MockYouTubeClientTestSuite{}
	suite.SetupTest()
	ctx := context.Background()
	video := &youtube.Video{ID: "test"}
	format := &youtube.Format{ItagNo: 999}
	expectedError := errors.New("stream not available")

	suite.mockClient.GetStreamContextFunc = func(ctx context.Context, video *youtube.Video, format *youtube.Format) (io.ReadCloser, int64, error) {
		return nil, 0, expectedError
	}

	// Act
	stream, size, err := suite.mockClient.GetStreamContext(ctx, video, format)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, stream)
	assert.Equal(t, int64(0), size)
	assert.Equal(t, expectedError, err)
}

// TestMockYouTubeClient_SetMaxRoutines_Normal は SetMaxRoutines メソッドの正常系テスト
func TestMockYouTubeClient_SetMaxRoutines_Normal(t *testing.T) {
	// Arrange
	suite := &MockYouTubeClientTestSuite{}
	suite.SetupTest()
	maxRoutines := 10
	called := false

	suite.mockClient.SetMaxRoutinesFunc = func(maxRoutines int) {
		called = true
	}

	// Act
	suite.mockClient.SetMaxRoutines(maxRoutines)

	// Assert
	assert.True(t, called)
}

// TestMockYouTubeClient_SetChunkSize_Normal は SetChunkSize メソッドの正常系テスト
func TestMockYouTubeClient_SetChunkSize_Normal(t *testing.T) {
	// Arrange
	suite := &MockYouTubeClientTestSuite{}
	suite.SetupTest()
	chunkSize := int64(1024)
	called := false

	suite.mockClient.SetChunkSizeFunc = func(chunkSize int64) {
		called = true
	}

	// Act
	suite.mockClient.SetChunkSize(chunkSize)

	// Assert
	assert.True(t, called)
}

// TestMockYouTubeClient_AllMethods_Nil は全てのメソッドのnilテスト
func TestMockYouTubeClient_AllMethods_Nil(t *testing.T) {
	// Arrange
	suite := &MockYouTubeClientTestSuite{}
	suite.SetupTest()
	ctx := context.Background()

	// Act & Assert
	// GetVideoContext
	result1, err1 := suite.mockClient.GetVideoContext(ctx, "test")
	assert.NoError(t, err1)
	assert.Nil(t, result1)

	// GetPlaylistContext
	result2, err2 := suite.mockClient.GetPlaylistContext(ctx, "test")
	assert.NoError(t, err2)
	assert.Nil(t, result2)

	// VideoFromPlaylistEntryContext
	result3, err3 := suite.mockClient.VideoFromPlaylistEntryContext(ctx, &youtube.PlaylistEntry{})
	assert.NoError(t, err3)
	assert.Nil(t, result3)

	// GetStreamContext
	stream, size, err4 := suite.mockClient.GetStreamContext(ctx, &youtube.Video{}, &youtube.Format{})
	assert.NoError(t, err4)
	assert.Nil(t, stream)
	assert.Equal(t, int64(0), size)

	// SetMaxRoutines (should not panic)
	suite.mockClient.SetMaxRoutines(10)

	// SetChunkSize (should not panic)
	suite.mockClient.SetChunkSize(1024)
}

// #==============================================================#
// ##         MockFileSystem Tests                              ##
// #==============================================================#
type MockFileSystemTestSuite struct {
	mockFileSystem *MockFileSystem
}

func (suite *MockFileSystemTestSuite) SetupTest() {
	suite.mockFileSystem = &MockFileSystem{}
}

// TestMockFileSystem_MkdirAll_Normal は MkdirAll メソッドの正常系テスト
func TestMockFileSystem_MkdirAll_Normal(t *testing.T) {
	// Arrange
	suite := &MockFileSystemTestSuite{}
	suite.SetupTest()
	path := "/test/path"
	perm := uint32(0755)

	suite.mockFileSystem.MkdirAllFunc = func(path string, perm uint32) error {
		return nil
	}

	// Act
	err := suite.mockFileSystem.MkdirAll(path, perm)

	// Assert
	assert.NoError(t, err)
}

// TestMockFileSystem_MkdirAll_Error は MkdirAll メソッドのエラー系テスト
func TestMockFileSystem_MkdirAll_Error(t *testing.T) {
	// Arrange
	suite := &MockFileSystemTestSuite{}
	suite.SetupTest()
	path := "/invalid/path"
	perm := uint32(0755)
	expectedError := errors.New("permission denied")

	suite.mockFileSystem.MkdirAllFunc = func(path string, perm uint32) error {
		return expectedError
	}

	// Act
	err := suite.mockFileSystem.MkdirAll(path, perm)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

// TestMockFileSystem_Create_Normal は Create メソッドの正常系テスト
func TestMockFileSystem_Create_Normal(t *testing.T) {
	// Arrange
	suite := &MockFileSystemTestSuite{}
	suite.SetupTest()
	name := "test.txt"
	expectedFile := &MockFile{}

	suite.mockFileSystem.CreateFunc = func(name string) (File, error) {
		return expectedFile, nil
	}

	// Act
	result, err := suite.mockFileSystem.Create(name)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedFile, result)
}

// TestMockFileSystem_Create_Error は Create メソッドのエラー系テスト
func TestMockFileSystem_Create_Error(t *testing.T) {
	// Arrange
	suite := &MockFileSystemTestSuite{}
	suite.SetupTest()
	name := "invalid.txt"
	expectedError := errors.New("file creation failed")

	suite.mockFileSystem.CreateFunc = func(name string) (File, error) {
		return nil, expectedError
	}

	// Act
	result, err := suite.mockFileSystem.Create(name)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

// TestMockFileSystem_Remove_Normal は Remove メソッドの正常系テスト
func TestMockFileSystem_Remove_Normal(t *testing.T) {
	// Arrange
	suite := &MockFileSystemTestSuite{}
	suite.SetupTest()
	name := "test.txt"

	suite.mockFileSystem.RemoveFunc = func(name string) error {
		return nil
	}

	// Act
	err := suite.mockFileSystem.Remove(name)

	// Assert
	assert.NoError(t, err)
}

// TestMockFileSystem_Remove_Error は Remove メソッドのエラー系テスト
func TestMockFileSystem_Remove_Error(t *testing.T) {
	// Arrange
	suite := &MockFileSystemTestSuite{}
	suite.SetupTest()
	name := "nonexistent.txt"
	expectedError := errors.New("file not found")

	suite.mockFileSystem.RemoveFunc = func(name string) error {
		return expectedError
	}

	// Act
	err := suite.mockFileSystem.Remove(name)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

// TestMockFileSystem_Stat_Normal は Stat メソッドの正常系テスト
func TestMockFileSystem_Stat_Normal(t *testing.T) {
	// Arrange
	suite := &MockFileSystemTestSuite{}
	suite.SetupTest()
	name := "test.txt"
	expectedFileInfo := &MockFileInfo{}

	suite.mockFileSystem.StatFunc = func(name string) (FileInfo, error) {
		return expectedFileInfo, nil
	}

	// Act
	result, err := suite.mockFileSystem.Stat(name)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedFileInfo, result)
}

// TestMockFileSystem_Stat_Error は Stat メソッドのエラー系テスト
func TestMockFileSystem_Stat_Error(t *testing.T) {
	// Arrange
	suite := &MockFileSystemTestSuite{}
	suite.SetupTest()
	name := "nonexistent.txt"
	expectedError := errors.New("file not found")

	suite.mockFileSystem.StatFunc = func(name string) (FileInfo, error) {
		return nil, expectedError
	}

	// Act
	result, err := suite.mockFileSystem.Stat(name)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

// TestMockFileSystem_AllMethods_Nil は全てのメソッドのnilテスト
func TestMockFileSystem_AllMethods_Nil(t *testing.T) {
	// Arrange
	suite := &MockFileSystemTestSuite{}
	suite.SetupTest()

	// Act & Assert
	// MkdirAll
	err1 := suite.mockFileSystem.MkdirAll("/test", 0755)
	assert.NoError(t, err1)

	// Create
	result2, err2 := suite.mockFileSystem.Create("test.txt")
	assert.NoError(t, err2)
	assert.NotNil(t, result2)

	// Remove
	err3 := suite.mockFileSystem.Remove("test.txt")
	assert.NoError(t, err3)

	// Stat
	result4, err4 := suite.mockFileSystem.Stat("test.txt")
	assert.NoError(t, err4)
	assert.NotNil(t, result4)
}

// #==============================================================#
// ##         MockFile Tests                                    ##
// #==============================================================#
type MockFileTestSuite struct {
	mockFile *MockFile
}

func (suite *MockFileTestSuite) SetupTest() {
	suite.mockFile = &MockFile{}
}

// TestMockFile_Write_Normal は Write メソッドの正常系テスト
func TestMockFile_Write_Normal(t *testing.T) {
	// Arrange
	suite := &MockFileTestSuite{}
	suite.SetupTest()
	data := []byte("test data")
	expectedN := len(data)

	suite.mockFile.WriteFunc = func(p []byte) (n int, err error) {
		return len(p), nil
	}

	// Act
	n, err := suite.mockFile.Write(data)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedN, n)
}

// TestMockFile_Write_Error は Write メソッドのエラー系テスト
func TestMockFile_Write_Error(t *testing.T) {
	// Arrange
	suite := &MockFileTestSuite{}
	suite.SetupTest()
	data := []byte("test data")
	expectedError := errors.New("write failed")

	suite.mockFile.WriteFunc = func(p []byte) (n int, err error) {
		return 0, expectedError
	}

	// Act
	n, err := suite.mockFile.Write(data)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, expectedError, err)
}

// TestMockFile_Close_Normal は Close メソッドの正常系テスト
func TestMockFile_Close_Normal(t *testing.T) {
	// Arrange
	suite := &MockFileTestSuite{}
	suite.SetupTest()

	suite.mockFile.CloseFunc = func() error {
		return nil
	}

	// Act
	err := suite.mockFile.Close()

	// Assert
	assert.NoError(t, err)
}

// TestMockFile_Close_Error は Close メソッドのエラー系テスト
func TestMockFile_Close_Error(t *testing.T) {
	// Arrange
	suite := &MockFileTestSuite{}
	suite.SetupTest()
	expectedError := errors.New("close failed")

	suite.mockFile.CloseFunc = func() error {
		return expectedError
	}

	// Act
	err := suite.mockFile.Close()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

// TestMockFile_Stat_Normal は Stat メソッドの正常系テスト
func TestMockFile_Stat_Normal(t *testing.T) {
	// Arrange
	suite := &MockFileTestSuite{}
	suite.SetupTest()
	expectedFileInfo := &MockFileInfo{}

	suite.mockFile.StatFunc = func() (FileInfo, error) {
		return expectedFileInfo, nil
	}

	// Act
	result, err := suite.mockFile.Stat()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedFileInfo, result)
}

// TestMockFile_Stat_Error は Stat メソッドのエラー系テスト
func TestMockFile_Stat_Error(t *testing.T) {
	// Arrange
	suite := &MockFileTestSuite{}
	suite.SetupTest()
	expectedError := errors.New("stat failed")

	suite.mockFile.StatFunc = func() (FileInfo, error) {
		return nil, expectedError
	}

	// Act
	result, err := suite.mockFile.Stat()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

// TestMockFile_AllMethods_Nil は全てのメソッドのnilテスト
func TestMockFile_AllMethods_Nil(t *testing.T) {
	// Arrange
	suite := &MockFileTestSuite{}
	suite.SetupTest()

	// Act & Assert
	// Write
	n, err1 := suite.mockFile.Write([]byte("test"))
	assert.NoError(t, err1)
	assert.Equal(t, 4, n)

	// Close
	err2 := suite.mockFile.Close()
	assert.NoError(t, err2)

	// Stat
	result, err3 := suite.mockFile.Stat()
	assert.NoError(t, err3)
	assert.NotNil(t, result)
}

// #==============================================================#
// ##         MockFileInfo Tests                                ##
// #==============================================================#
type MockFileInfoTestSuite struct {
	mockFileInfo *MockFileInfo
}

func (suite *MockFileInfoTestSuite) SetupTest() {
	suite.mockFileInfo = &MockFileInfo{}
}

// TestMockFileInfo_Size_Normal は Size メソッドの正常系テスト
func TestMockFileInfo_Size_Normal(t *testing.T) {
	// Arrange
	suite := &MockFileInfoTestSuite{}
	suite.SetupTest()
	expectedSize := int64(1024)

	suite.mockFileInfo.SizeFunc = func() int64 {
		return expectedSize
	}

	// Act
	result := suite.mockFileInfo.Size()

	// Assert
	assert.Equal(t, expectedSize, result)
}

// TestMockFileInfo_Size_Nil は Size メソッドのnilテスト
func TestMockFileInfo_Size_Nil(t *testing.T) {
	// Arrange
	suite := &MockFileInfoTestSuite{}
	suite.SetupTest()

	// Act
	result := suite.mockFileInfo.Size()

	// Assert
	assert.Equal(t, int64(0), result)
}

// #==============================================================#
// ##         MockVideoProcessor Tests                          ##
// #==============================================================#
type MockVideoProcessorTestSuite struct {
	mockProcessor *MockVideoProcessor
}

func (suite *MockVideoProcessorTestSuite) SetupTest() {
	suite.mockProcessor = &MockVideoProcessor{}
}

// TestMockVideoProcessor_MergeVideoAudio_Normal は MergeVideoAudio メソッドの正常系テスト
func TestMockVideoProcessor_MergeVideoAudio_Normal(t *testing.T) {
	// Arrange
	suite := &MockVideoProcessorTestSuite{}
	suite.SetupTest()
	videoPath := "video.mp4"
	audioPath := "audio.mp3"
	outputPath := "output.mp4"

	suite.mockProcessor.MergeVideoAudioFunc = func(videoPath, audioPath, outputPath string) error {
		return nil
	}

	// Act
	err := suite.mockProcessor.MergeVideoAudio(videoPath, audioPath, outputPath)

	// Assert
	assert.NoError(t, err)
}

// TestMockVideoProcessor_MergeVideoAudio_Error は MergeVideoAudio メソッドのエラー系テスト
func TestMockVideoProcessor_MergeVideoAudio_Error(t *testing.T) {
	// Arrange
	suite := &MockVideoProcessorTestSuite{}
	suite.SetupTest()
	videoPath := "invalid_video.mp4"
	audioPath := "invalid_audio.mp3"
	outputPath := "output.mp4"
	expectedError := errors.New("merge failed")

	suite.mockProcessor.MergeVideoAudioFunc = func(videoPath, audioPath, outputPath string) error {
		return expectedError
	}

	// Act
	err := suite.mockProcessor.MergeVideoAudio(videoPath, audioPath, outputPath)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

// TestMockVideoProcessor_MergeVideoAudio_Nil は MergeVideoAudio メソッドのnilテスト
func TestMockVideoProcessor_MergeVideoAudio_Nil(t *testing.T) {
	// Arrange
	suite := &MockVideoProcessorTestSuite{}
	suite.SetupTest()

	// Act
	err := suite.mockProcessor.MergeVideoAudio("video.mp4", "audio.mp3", "output.mp4")

	// Assert
	assert.NoError(t, err)
}

// #==============================================================#
// ##         MockTimeProvider Tests                            ##
// #==============================================================#
type MockTimeProviderTestSuite struct {
	mockProvider *MockTimeProvider
}

func (suite *MockTimeProviderTestSuite) SetupTest() {
	suite.mockProvider = &MockTimeProvider{}
}

// TestMockTimeProvider_Now_Normal は Now メソッドの正常系テスト
func TestMockTimeProvider_Now_Normal(t *testing.T) {
	// Arrange
	suite := &MockTimeProviderTestSuite{}
	suite.SetupTest()
	expectedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	suite.mockProvider.NowFunc = func() time.Time {
		return expectedTime
	}

	// Act
	result := suite.mockProvider.Now()

	// Assert
	assert.Equal(t, expectedTime, result)
}

// TestMockTimeProvider_Now_Nil は Now メソッドのnilテスト
func TestMockTimeProvider_Now_Nil(t *testing.T) {
	// Arrange
	suite := &MockTimeProviderTestSuite{}
	suite.SetupTest()

	// Act
	result := suite.mockProvider.Now()

	// Assert
	assert.NotZero(t, result)
}

// TestMockTimeProvider_LoadLocation_Normal は LoadLocation メソッドの正常系テスト
func TestMockTimeProvider_LoadLocation_Normal(t *testing.T) {
	// Arrange
	suite := &MockTimeProviderTestSuite{}
	suite.SetupTest()
	name := "Asia/Tokyo"
	expectedLocation, _ := time.LoadLocation(name)

	suite.mockProvider.LoadLocationFunc = func(name string) (*time.Location, error) {
		return expectedLocation, nil
	}

	// Act
	result, err := suite.mockProvider.LoadLocation(name)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedLocation, result)
}

// TestMockTimeProvider_LoadLocation_Error は LoadLocation メソッドのエラー系テスト
func TestMockTimeProvider_LoadLocation_Error(t *testing.T) {
	// Arrange
	suite := &MockTimeProviderTestSuite{}
	suite.SetupTest()
	name := "Invalid/Timezone"
	expectedError := errors.New("unknown time zone")

	suite.mockProvider.LoadLocationFunc = func(name string) (*time.Location, error) {
		return nil, expectedError
	}

	// Act
	result, err := suite.mockProvider.LoadLocation(name)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

// TestMockTimeProvider_LoadLocation_Nil は LoadLocation メソッドのnilテスト
func TestMockTimeProvider_LoadLocation_Nil(t *testing.T) {
	// Arrange
	suite := &MockTimeProviderTestSuite{}
	suite.SetupTest()
	name := "UTC"

	// Act
	result, err := suite.mockProvider.LoadLocation(name)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "UTC", result.String())
}
