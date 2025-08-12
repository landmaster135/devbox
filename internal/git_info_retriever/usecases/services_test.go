package usecases

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#

// MockGitHubService はGitHubServiceのモック
type MockGitHubService struct {
	mock.Mock
}

func (m *MockGitHubService) CreateGitHubClient(ctx context.Context, token string) GitHubClient {
	args := m.Called(ctx, token)
	return args.Get(0).(GitHubClient)
}

func (m *MockGitHubService) GetRepoInfo(ctx context.Context, client GitHubClient, isThreading bool, username string) ([]RepoInfo, error) {
	args := m.Called(ctx, client, isThreading, username)
	return args.Get(0).([]RepoInfo), args.Error(1)
}

// MockGitHubClient はGitHubClientのモック
type MockGitHubClient struct {
	mock.Mock
}

func (m *MockGitHubClient) ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	args := m.Called(ctx, user, opts)
	return args.Get(0).([]*github.Repository), args.Get(1).(*github.Response), args.Error(2)
}

func (m *MockGitHubClient) ListRepoLanguages(ctx context.Context, owner string, repo string) (map[string]int, *github.Response, error) {
	args := m.Called(ctx, owner, repo)
	return args.Get(0).(map[string]int), args.Get(1).(*github.Response), args.Error(2)
}

func (m *MockGitHubClient) ListPullRequests(ctx context.Context, user string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
	args := m.Called(ctx, user, repo, opts)
	return args.Get(0).([]*github.PullRequest), args.Get(1).(*github.Response), args.Error(2)
}

func (m *MockGitHubClient) GetUser(ctx context.Context, user string) (*github.User, *github.Response, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*github.User), args.Get(1).(*github.Response), args.Error(2)
}

// MockFileWriter はFileWriterのモック
type MockFileWriter struct {
	mock.Mock
}

func (m *MockFileWriter) WriteToFile(filePath, content string) error {
	args := m.Called(filePath, content)
	return args.Error(0)
}

func (m *MockFileWriter) EnsureDirectory(dirPath string) error {
	args := m.Called(dirPath)
	return args.Error(0)
}

// #==============================================================#
// ##          Test Helper Functions                             ##
// #==============================================================#

// createTestRepoInfo はテスト用のRepoInfoを作成する
func createTestRepoInfo() []RepoInfo {
	return []RepoInfo{
		{
			Name:             "test-repo-1",
			Description:      "Test repository 1",
			IsPrivate:        false,
			HttpUrl:          "https://github.com/testuser/test-repo-1",
			Language:         "Go",
			Languages:        map[string]int{"Go": 1000, "JavaScript": 500},
			RepoCreatedAt:    "2023-01-01T00:00:00Z",
			RepoUpdatedAt:    "2023-12-01T00:00:00Z",
			StargazersCount:  10,
			ForksCount:       5,
			IssuesCount:      2,
			PullsCount:       3,
			Size:             1024,
			SubscribersCount: 8,
			IsArchived:       false,
		},
		{
			Name:             "test-repo-2",
			Description:      "Test repository 2",
			IsPrivate:        true,
			HttpUrl:          "https://github.com/testuser/test-repo-2",
			Language:         "Python",
			Languages:        map[string]int{"Python": 2000, "HTML": 300},
			RepoCreatedAt:    "2023-02-01T00:00:00Z",
			RepoUpdatedAt:    "2023-11-01T00:00:00Z",
			StargazersCount:  20,
			ForksCount:       8,
			IssuesCount:      1,
			PullsCount:       5,
			Size:             2048,
			SubscribersCount: 15,
			IsArchived:       false,
		},
	}
}

// createTestGitHubUser はテスト用のGitHubユーザーを作成する
func createTestGitHubUser() *github.User {
	login := "testuser"
	return &github.User{
		Login: &login,
	}
}

// setupServiceWithMocks はモックを使用してServiceをセットアップする
func setupServiceWithMocks(t *testing.T) (*Service, *MockGitHubService, *MockFileWriter) {
	mockGitHubService := new(MockGitHubService)
	mockFileWriter := new(MockFileWriter)

	service := &Service{
		githubService: mockGitHubService,
		fileWriter:    mockFileWriter,
	}

	return service, mockGitHubService, mockFileWriter
}

// #==============================================================#
// ##          Service Tests                                     ##
// #==============================================================#

func TestService_RetrieveRepositoryInfo_Normal(t *testing.T) {
	const (
		testService  = "github"
		testToken    = "test-token"
		testFilePath = "/tmp/test-output.json"
		testUsername = "testuser"
	)

	tests := []struct {
		name                    string
		service                 string
		token                   string
		filePath                string
		setupGitHubServiceMock  func(*MockGitHubService, *MockGitHubClient)
		setupFileWriterMock     func(*MockFileWriter)
		expectError             bool
		expectedErrorMessage    string
		expectedResultContains  string
	}{
		{
			name:     "WithFilePath_Normal",
			service:  testService,
			token:    testToken,
			filePath: testFilePath,
			setupGitHubServiceMock: func(mockGitHubService *MockGitHubService, mockGitHubClient *MockGitHubClient) {
				mockGitHubService.On("CreateGitHubClient", mock.Anything, testToken).Return(mockGitHubClient)
				mockGitHubClient.On("GetUser", mock.Anything, "").Return(createTestGitHubUser(), &github.Response{}, nil)
				mockGitHubService.On("GetRepoInfo", mock.Anything, mockGitHubClient, true, testUsername).Return(createTestRepoInfo(), nil)
			},
			setupFileWriterMock: func(mockFileWriter *MockFileWriter) {
				mockFileWriter.On("WriteToFile", testFilePath, mock.AnythingOfType("string")).Return(nil)
			},
			expectError:            false,
			expectedResultContains: "test-repo-1",
		},
		{
			name:     "WithoutFilePath_Normal",
			service:  testService,
			token:    testToken,
			filePath: "",
			setupGitHubServiceMock: func(mockGitHubService *MockGitHubService, mockGitHubClient *MockGitHubClient) {
				mockGitHubService.On("CreateGitHubClient", mock.Anything, testToken).Return(mockGitHubClient)
				mockGitHubClient.On("GetUser", mock.Anything, "").Return(createTestGitHubUser(), &github.Response{}, nil)
				mockGitHubService.On("GetRepoInfo", mock.Anything, mockGitHubClient, true, testUsername).Return(createTestRepoInfo(), nil)
			},
			setupFileWriterMock:    func(mockFileWriter *MockFileWriter) {},
			expectError:            false,
			expectedResultContains: "test-repo-1",
		},
		{
			name:     "UnsupportedService_Error",
			service:  "gitlab",
			token:    testToken,
			filePath: "",
			setupGitHubServiceMock: func(mockGitHubService *MockGitHubService, mockGitHubClient *MockGitHubClient) {},
			setupFileWriterMock:    func(mockFileWriter *MockFileWriter) {},
			expectError:            true,
			expectedErrorMessage:   "サポートされていないサービスです: gitlab",
		},
		{
			name:     "GitHubUserError_Error",
			service:  testService,
			token:    testToken,
			filePath: "",
			setupGitHubServiceMock: func(mockGitHubService *MockGitHubService, mockGitHubClient *MockGitHubClient) {
				mockGitHubService.On("CreateGitHubClient", mock.Anything, testToken).Return(mockGitHubClient)
				mockGitHubClient.On("GetUser", mock.Anything, "").Return((*github.User)(nil), &github.Response{}, errors.New("user not found"))
			},
			setupFileWriterMock:  func(mockFileWriter *MockFileWriter) {},
			expectError:          true,
			expectedErrorMessage: "ユーザー情報の取得に失敗しました",
		},
		{
			name:     "GitHubRepoInfoError_Error",
			service:  testService,
			token:    testToken,
			filePath: "",
			setupGitHubServiceMock: func(mockGitHubService *MockGitHubService, mockGitHubClient *MockGitHubClient) {
				mockGitHubService.On("CreateGitHubClient", mock.Anything, testToken).Return(mockGitHubClient)
				mockGitHubClient.On("GetUser", mock.Anything, "").Return(createTestGitHubUser(), &github.Response{}, nil)
				mockGitHubService.On("GetRepoInfo", mock.Anything, mockGitHubClient, true, testUsername).Return([]RepoInfo{}, errors.New("repo info error"))
			},
			setupFileWriterMock:  func(mockFileWriter *MockFileWriter) {},
			expectError:          true,
			expectedErrorMessage: "リポジトリ情報の取得に失敗しました",
		},
		{
			name:     "FileWriteError_Error",
			service:  testService,
			token:    testToken,
			filePath: testFilePath,
			setupGitHubServiceMock: func(mockGitHubService *MockGitHubService, mockGitHubClient *MockGitHubClient) {
				mockGitHubService.On("CreateGitHubClient", mock.Anything, testToken).Return(mockGitHubClient)
				mockGitHubClient.On("GetUser", mock.Anything, "").Return(createTestGitHubUser(), &github.Response{}, nil)
				mockGitHubService.On("GetRepoInfo", mock.Anything, mockGitHubClient, true, testUsername).Return(createTestRepoInfo(), nil)
			},
			setupFileWriterMock: func(mockFileWriter *MockFileWriter) {
				mockFileWriter.On("WriteToFile", testFilePath, mock.AnythingOfType("string")).Return(errors.New("write error"))
			},
			expectError:          true,
			expectedErrorMessage: "ファイル保存に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockGitHubService, mockFileWriter := setupServiceWithMocks(t)
			mockGitHubClient := new(MockGitHubClient)

			tt.setupGitHubServiceMock(mockGitHubService, mockGitHubClient)
			tt.setupFileWriterMock(mockFileWriter)

			result, err := service.RetrieveRepositoryInfo(tt.service, tt.token, tt.filePath)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMessage)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, result, tt.expectedResultContains)
			}

			mockGitHubService.AssertExpectations(t)
			mockFileWriter.AssertExpectations(t)
			mockGitHubClient.AssertExpectations(t)
		})
	}
}

func TestNewService_Normal(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)

	// 型アサーションでServiceの実装を確認
	serviceImpl, ok := service.(*Service)
	assert.True(t, ok)
	assert.NotNil(t, serviceImpl.githubService)
	assert.NotNil(t, serviceImpl.fileWriter)
}

func TestNewServiceWithDependencies_Normal(t *testing.T) {
	mockGitHubService := new(MockGitHubService)
	mockFileWriter := new(MockFileWriter)

	service := NewServiceWithDependencies(mockGitHubService, mockFileWriter)
	assert.NotNil(t, service)

	// 型アサーションでServiceの実装を確認
	serviceImpl, ok := service.(*Service)
	assert.True(t, ok)
	assert.Equal(t, mockGitHubService, serviceImpl.githubService)
	assert.Equal(t, mockFileWriter, serviceImpl.fileWriter)
}

// #==============================================================#
// ##          GitHubServiceImpl Tests                           ##
// #==============================================================#

func TestGitHubServiceImpl_GetRepoInfo_Normal(t *testing.T) {
	const (
		testUsername = "testuser"
	)

	tests := []struct {
		name                string
		isThreading         bool
		setupMock           func(*MockGitHubClient)
		expectError         bool
		expectedErrorMsg    string
		expectedRepoCount   int
	}{
		{
			name:        "WithThreading_Normal",
			isThreading: true,
			setupMock: func(mockClient *MockGitHubClient) {
				// fetchRepositories用のモック
				repos := createTestGitHubRepositories()
				mockClient.On("ListRepositories", mock.Anything, testUsername, mock.AnythingOfType("*github.RepositoryListOptions")).Return(repos, &github.Response{NextPage: 0}, nil)

				// getRepoInfoInFormat用のモック（各リポジトリに対して）
				for _, repo := range repos {
					mockClient.On("ListPullRequests", mock.Anything, repo.GetOwner().GetLogin(), repo.GetName(), mock.AnythingOfType("*github.PullRequestListOptions")).Return([]*github.PullRequest{}, &github.Response{}, nil)
					mockClient.On("ListRepoLanguages", mock.Anything, repo.GetOwner().GetLogin(), repo.GetName()).Return(map[string]int{"Go": 1000}, &github.Response{}, nil)
				}
			},
			expectError:       false,
			expectedRepoCount: 2,
		},
		{
			name:        "WithoutThreading_Normal",
			isThreading: false,
			setupMock: func(mockClient *MockGitHubClient) {
				// fetchRepositories用のモック
				repos := createTestGitHubRepositories()
				mockClient.On("ListRepositories", mock.Anything, testUsername, mock.AnythingOfType("*github.RepositoryListOptions")).Return(repos, &github.Response{NextPage: 0}, nil)

				// getRepoInfoInFormat用のモック（各リポジトリに対して）
				for _, repo := range repos {
					mockClient.On("ListPullRequests", mock.Anything, repo.GetOwner().GetLogin(), repo.GetName(), mock.AnythingOfType("*github.PullRequestListOptions")).Return([]*github.PullRequest{}, &github.Response{}, nil)
					mockClient.On("ListRepoLanguages", mock.Anything, repo.GetOwner().GetLogin(), repo.GetName()).Return(map[string]int{"Go": 1000}, &github.Response{}, nil)
				}
			},
			expectError:       false,
			expectedRepoCount: 2,
		},
		{
			name:        "FetchRepositoriesError_Error",
			isThreading: false,
			setupMock: func(mockClient *MockGitHubClient) {
				mockClient.On("ListRepositories", mock.Anything, testUsername, mock.AnythingOfType("*github.RepositoryListOptions")).Return([]*github.Repository{}, &github.Response{}, errors.New("fetch error"))
			},
			expectError:      true,
			expectedErrorMsg: "fetch error",
		},
		{
			name:        "RateLimitError_Error",
			isThreading: false,
			setupMock: func(mockClient *MockGitHubClient) {
				rateLimitErr := &github.RateLimitError{
					Rate: github.Rate{
						Limit:     5000,
						Remaining: 0,
						Reset:     github.Timestamp{Time: time.Now().Add(time.Hour)},
					},
				}
				mockClient.On("ListRepositories", mock.Anything, testUsername, mock.AnythingOfType("*github.RepositoryListOptions")).Return([]*github.Repository{}, &github.Response{}, rateLimitErr)
			},
			expectError:      true,
			expectedErrorMsg: "GitHub APIのレート制限に達しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewGitHubService()
			mockClient := new(MockGitHubClient)
			tt.setupMock(mockClient)

			ctx := context.Background()
			result, err := service.GetRepoInfo(ctx, mockClient, tt.isThreading, testUsername)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedRepoCount)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestGitHubServiceImpl_fetchRepositories_Normal(t *testing.T) {
	const testUsername = "testuser"

	tests := []struct {
		name             string
		repoType         string
		setupMock        func(*MockGitHubClient)
		expectError      bool
		expectedErrorMsg string
		expectedCount    int
	}{
		{
			name:     "ValidRepoType_Normal",
			repoType: "owner",
			setupMock: func(mockClient *MockGitHubClient) {
				repos := createTestGitHubRepositories()
				mockClient.On("ListRepositories", mock.Anything, testUsername, mock.AnythingOfType("*github.RepositoryListOptions")).Return(repos, &github.Response{NextPage: 0}, nil)
			},
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:             "EmptyRepoType_Error",
			repoType:         "",
			setupMock:        func(mockClient *MockGitHubClient) {},
			expectError:      true,
			expectedErrorMsg: "リポジトリタイプが空です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &GitHubServiceImpl{}
			mockClient := new(MockGitHubClient)
			tt.setupMock(mockClient)

			ctx := context.Background()
			result, err := service.fetchRepositories(ctx, mockClient, tt.repoType, testUsername)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestGitHubServiceImpl_getRepoPulls_Normal(t *testing.T) {
	const (
		testOwner = "testowner"
		testRepo  = "testrepo"
	)

	tests := []struct {
		name             string
		state            string
		setupMock        func(*MockGitHubClient)
		expectError      bool
		expectedErrorMsg string
		expectedCount    int
	}{
		{
			name:  "ValidState_Normal",
			state: "all",
			setupMock: func(mockClient *MockGitHubClient) {
				pulls := createTestPullRequests()
				mockClient.On("ListPullRequests", mock.Anything, testOwner, testRepo, mock.AnythingOfType("*github.PullRequestListOptions")).Return(pulls, &github.Response{}, nil)
			},
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:             "EmptyState_Error",
			state:            "",
			setupMock:        func(mockClient *MockGitHubClient) {},
			expectError:      true,
			expectedErrorMsg: "stateが空です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &GitHubServiceImpl{}
			mockClient := new(MockGitHubClient)
			tt.setupMock(mockClient)

			repo := createTestGitHubRepository(testOwner, testRepo)
			ctx := context.Background()
			result, err := service.getRepoPulls(ctx, mockClient, repo, tt.state)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestGitHubServiceImpl_CreateGitHubClient_Normal(t *testing.T) {
	const testToken = "test-token"

	service := &GitHubServiceImpl{}
	ctx := context.Background()

	client := service.CreateGitHubClient(ctx, testToken)

	assert.NotNil(t, client)

	// 型アサーションでGitHubClientAdapterの実装を確認
	adapter, ok := client.(*GitHubClientAdapter)
	assert.True(t, ok)
	assert.NotNil(t, adapter.Client)
}

// #==============================================================#
// ##          FileWriterImpl Tests                              ##
// #==============================================================#

func TestFileWriterImpl_WriteToFile_Normal(t *testing.T) {
	const (
		testContent = "test content"
	)

	tests := []struct {
		name             string
		filePath         string
		content          string
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:        "ValidPath_Normal",
			filePath:    "/tmp/test-file.txt",
			content:     testContent,
			expectError: false,
		},
		{
			name:        "ValidPathWithDirectory_Normal",
			filePath:    "/tmp/test-dir/test-file.txt",
			content:     testContent,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewFileWriter()

			err := writer.WriteToFile(tt.filePath, tt.content)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileWriterImpl_EnsureDirectory_Normal(t *testing.T) {
	tests := []struct {
		name        string
		dirPath     string
		expectError bool
	}{
		{
			name:        "ValidDirectory_Normal",
			dirPath:     "/tmp/test-ensure-dir",
			expectError: false,
		},
		{
			name:        "CurrentDirectory_Normal",
			dirPath:     ".",
			expectError: false,
		},
		{
			name:        "EmptyDirectory_Normal",
			dirPath:     "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &FileWriterImpl{}

			err := writer.EnsureDirectory(tt.dirPath)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// #==============================================================#
// ##          Utility Functions for Tests                      ##
// #==============================================================#

// createTestGitHubRepositories はテスト用のGitHubリポジトリ配列を作成する
func createTestGitHubRepositories() []*github.Repository {
	return []*github.Repository{
		createTestGitHubRepository("testowner", "test-repo-1"),
		createTestGitHubRepository("testowner", "test-repo-2"),
	}
}

// createTestGitHubRepository はテスト用のGitHubリポジトリを作成する
func createTestGitHubRepository(owner, name string) *github.Repository {
	createdAt := github.Timestamp{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}
	updatedAt := github.Timestamp{Time: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)}

	return &github.Repository{
		Name:        &name,
		Description: github.String(fmt.Sprintf("Description for %s", name)),
		Private:     github.Bool(false),
		HTMLURL:     github.String(fmt.Sprintf("https://github.com/%s/%s", owner, name)),
		Language:    github.String("Go"),
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
		StargazersCount: github.Int(10),
		ForksCount:      github.Int(5),
		OpenIssuesCount: github.Int(2),
		Size:            github.Int(1024),
		SubscribersCount: github.Int(8),
		Archived:        github.Bool(false),
		Owner: &github.User{
			Login: &owner,
		},
	}
}

// createTestPullRequests はテスト用のプルリクエスト配列を作成する
func createTestPullRequests() []*github.PullRequest {
	return []*github.PullRequest{
		{
			ID:     github.Int64(1),
			Number: github.Int(1),
			Title:  github.String("Test PR 1"),
		},
		{
			ID:     github.Int64(2),
			Number: github.Int(2),
			Title:  github.String("Test PR 2"),
		},
	}
}

// formatGitHubTime のテスト
func TestFormatGitHubTime_Normal(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 12, 30, 45, 0, time.UTC)
	timestamp := github.Timestamp{Time: testTime}

	result := formatGitHubTime(timestamp)
	expected := "2023-01-01T12:30:45Z"

	assert.Equal(t, expected, result)
}

// #==============================================================#
// ##          GitHubClientAdapter Tests                         ##
// #==============================================================#

func TestGitHubClientAdapter_ListRepositories_Normal(t *testing.T) {
	const testUsername = "testuser"

	// 実際のGitHubクライアントを使用してアダプターをテスト
	service := &GitHubServiceImpl{}
	ctx := context.Background()
	client := service.CreateGitHubClient(ctx, "dummy-token")

	// アダプターの型アサーション
	adapter, ok := client.(*GitHubClientAdapter)
	assert.True(t, ok)
	assert.NotNil(t, adapter.Client)

	// メソッドが存在することを確認（実際のAPI呼び出しはしない）
	opts := &github.RepositoryListOptions{
		Type: "owner",
		ListOptions: github.ListOptions{PerPage: 1},
	}

	// ダミートークンなのでエラーになるが、メソッドが呼び出せることを確認
	_, _, err := adapter.ListRepositories(ctx, testUsername, opts)
	assert.Error(t, err) // 認証エラーが期待される
}

func TestGitHubClientAdapter_ListRepoLanguages_Normal(t *testing.T) {
	const (
		testOwner = "testowner"
		testRepo  = "testrepo"
	)

	service := &GitHubServiceImpl{}
	ctx := context.Background()
	client := service.CreateGitHubClient(ctx, "dummy-token")

	adapter, ok := client.(*GitHubClientAdapter)
	assert.True(t, ok)

	// メソッドが存在することを確認
	_, _, err := adapter.ListRepoLanguages(ctx, testOwner, testRepo)
	assert.Error(t, err) // 認証エラーが期待される
}

func TestGitHubClientAdapter_ListPullRequests_Normal(t *testing.T) {
	const (
		testUser = "testuser"
		testRepo = "testrepo"
	)

	service := &GitHubServiceImpl{}
	ctx := context.Background()
	client := service.CreateGitHubClient(ctx, "dummy-token")

	adapter, ok := client.(*GitHubClientAdapter)
	assert.True(t, ok)

	opts := &github.PullRequestListOptions{State: "all"}

	// メソッドが存在することを確認
	_, _, err := adapter.ListPullRequests(ctx, testUser, testRepo, opts)
	assert.Error(t, err) // 認証エラーが期待される
}

func TestGitHubClientAdapter_GetUser_Normal(t *testing.T) {
	const testUser = "testuser"

	service := &GitHubServiceImpl{}
	ctx := context.Background()
	client := service.CreateGitHubClient(ctx, "dummy-token")

	adapter, ok := client.(*GitHubClientAdapter)
	assert.True(t, ok)

	// メソッドが存在することを確認
	_, _, err := adapter.GetUser(ctx, testUser)
	assert.Error(t, err) // 認証エラーが期待される
}

// #==============================================================#
// ##          Additional FileWriter Tests                      ##
// #==============================================================#

func TestFileWriterImpl_WriteToFile_ErrorCases_Normal(t *testing.T) {
	tests := []struct {
		name             string
		filePath         string
		content          string
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:             "InvalidPath_Error",
			filePath:         "/root/no-permission/test.txt",
			content:          "test content",
			expectError:      true,
			expectedErrorMsg: "ディレクトリの作成に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &FileWriterImpl{}

			err := writer.WriteToFile(tt.filePath, tt.content)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileWriterImpl_EnsureDirectory_ErrorCases_Normal(t *testing.T) {
	tests := []struct {
		name             string
		dirPath          string
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:             "InvalidPermission_Error",
			dirPath:          "/root/no-permission-dir",
			expectError:      true,
			expectedErrorMsg: "ディレクトリの作成に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &FileWriterImpl{}

			err := writer.EnsureDirectory(tt.dirPath)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// #==============================================================#
// ##          Additional GitHub Service Tests                  ##
// #==============================================================#

func TestGitHubServiceImpl_getRepoInfoInFormat_ErrorCases_Normal(t *testing.T) {
	const (
		testOwner = "testowner"
		testRepo  = "testrepo"
	)

	tests := []struct {
		name             string
		setupMock        func(*MockGitHubClient)
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "PullRequestsError_Error",
			setupMock: func(mockClient *MockGitHubClient) {
				mockClient.On("ListPullRequests", mock.Anything, testOwner, testRepo, mock.AnythingOfType("*github.PullRequestListOptions")).Return([]*github.PullRequest{}, &github.Response{}, errors.New("pull requests error"))
			},
			expectError:      true,
			expectedErrorMsg: "pull requests error",
		},
		{
			name: "LanguagesError_Error",
			setupMock: func(mockClient *MockGitHubClient) {
				mockClient.On("ListPullRequests", mock.Anything, testOwner, testRepo, mock.AnythingOfType("*github.PullRequestListOptions")).Return([]*github.PullRequest{}, &github.Response{}, nil)
				mockClient.On("ListRepoLanguages", mock.Anything, testOwner, testRepo).Return(map[string]int{}, &github.Response{}, errors.New("languages error"))
			},
			expectError:      true,
			expectedErrorMsg: "languages error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &GitHubServiceImpl{}
			mockClient := new(MockGitHubClient)
			tt.setupMock(mockClient)

			repo := createTestGitHubRepository(testOwner, testRepo)
			ctx := context.Background()
			result, err := service.getRepoInfoInFormat(ctx, mockClient, repo)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
