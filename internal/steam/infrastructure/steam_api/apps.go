package steam_api

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// AppsService はSteam Apps APIのサービス
type AppsService struct {
	client *Client
}

// NewAppsService は新しいAppsServiceを作成します
func NewAppsService(client *Client) *AppsService {
	return &AppsService{
		client: client,
	}
}

// GetAppDetails はアプリケーションの詳細情報を取得します
func (a *AppsService) GetAppDetails(ctx context.Context, appID int, country, filters string) (*AppDetails, error) {
	if !validateAppID(appID) {
		return nil, fmt.Errorf("invalid App ID: %d", appID)
	}

	if country == "" {
		country = "US"
	}

	if filters == "" {
		filters = "basic"
	}

	params := map[string]any{
		"appids":  appID,
		"cc":      country,
		"filters": filters,
	}

	result, err := a.client.RequestWithoutKey(ctx, "GET", APIAppDetailsURL, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get app details: %w", err)
	}

	// レスポンスをパース
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var response AppDetailsResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	appIDStr := strconv.Itoa(appID)
	appData, exists := response[appIDStr]
	if !exists {
		return nil, fmt.Errorf("app data not found for ID: %d", appID)
	}

	if !appData.Success {
		return nil, fmt.Errorf("failed to get app details for ID: %d", appID)
	}

	if appData.Data == nil {
		return nil, fmt.Errorf("app data is nil for ID: %d", appID)
	}

	return appData.Data, nil
}

// GetUserStats はユーザーの特定アプリの統計情報を取得します
func (a *AppsService) GetUserStats(ctx context.Context, steamID string, appID int) (*UserStats, error) {
	if !validateSteamID(steamID) {
		return nil, fmt.Errorf("invalid Steam ID: %s", steamID)
	}

	if !validateAppID(appID) {
		return nil, fmt.Errorf("invalid App ID: %d", appID)
	}

	params := map[string]any{
		"steamid": steamID,
		"appid":   appID,
	}

	result, err := a.client.Request(ctx, "GET", EndpointUserStatsForGame, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	// レスポンスをパース
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var response UserStatsResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &UserStats{
		SteamID:      response.PlayerstatsResponse.SteamID,
		GameName:     response.PlayerstatsResponse.GameName,
		Stats:        response.PlayerstatsResponse.Stats,
		Achievements: response.PlayerstatsResponse.Achievements,
	}, nil
}

// GetUserAchievements はユーザーの特定アプリの実績情報を取得します
func (a *AppsService) GetUserAchievements(ctx context.Context, steamID string, appID int, language string) ([]Achievement, error) {
	if !validateSteamID(steamID) {
		return nil, fmt.Errorf("invalid Steam ID: %s", steamID)
	}

	if !validateAppID(appID) {
		return nil, fmt.Errorf("invalid App ID: %d", appID)
	}

	if language == "" {
		language = "en"
	}

	params := map[string]any{
		"steamid": steamID,
		"appid":   appID,
		"l":       language,
	}

	result, err := a.client.Request(ctx, "GET", EndpointUserPlayerAchievements, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get user achievements: %w", err)
	}

	// レスポンスをパース
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var response PlayerAchievementsResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.PlayerstatsResponse.Success {
		return nil, fmt.Errorf("failed to get achievements for user %s and app %d", steamID, appID)
	}

	return response.PlayerstatsResponse.Achievements, nil
}

// SearchGames はゲームを検索します（Store API使用）
func (a *AppsService) SearchGames(ctx context.Context, term, country string) ([]SearchResult, error) {
	if term == "" {
		return nil, fmt.Errorf("search term is required")
	}

	if country == "" {
		country = "US"
	}

	params := map[string]any{
		"f":     "games",
		"cc":    country,
		"realm": 1,
		"l":     "english",
	}

	url := buildURLWithParamsForSearch(APIAppSearchURL, term, params)

	result, err := a.client.RequestWithoutKey(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search games: %w", err)
	}

	// Store APIの検索結果はHTMLで返されるため、簡単な解析を行う
	// 実際の実装では、HTMLパーサーを使用してより詳細な解析を行う必要がある
	resultStr, ok := result.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	// 簡単な実装として、空のスライスを返す
	// 実際の実装では、HTMLを解析してSearchResult構造体のスライスを返す
	_ = resultStr
	return []SearchResult{}, nil
}

// SearchResult はゲーム検索結果の構造体
type SearchResult struct {
	ID    []int  `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
	Image string `json:"img"`
	Link  string `json:"link"`
}

// GetPublishedFileDetails はWorkshopファイルの詳細を取得します
func (a *AppsService) GetPublishedFileDetails(ctx context.Context, publishedFileIDs []int, options *PublishedFileOptions) (any, error) {
	if len(publishedFileIDs) == 0 {
		return nil, fmt.Errorf("published file IDs are required")
	}

	if options == nil {
		options = &PublishedFileOptions{}
	}

	params := map[string]any{
		"includetags":               options.IncludeTags,
		"includeadditionalpreviews": options.IncludeAdditionalPreviews,
		"includechildren":           options.IncludeChildren,
		"includekvtags":             options.IncludeKVTags,
		"includevotes":              options.IncludeVotes,
		"short_description":         options.ShortDescription,
		"includeforsaledata":        options.IncludeForSaleData,
		"includemetadata":           options.IncludeMetadata,
		"language":                  options.Language,
		"return_playtime_stats":     options.ReturnPlaytimeStats,
		"appid":                     options.AppID,
		"strip_description_bbcode":  options.StripDescriptionBBCode,
		"desired_revision":          options.DesiredRevision,
		"includereactions":          options.IncludeReactions,
		"admin_query":               options.AdminQuery,
	}

	// Published file IDsをパラメータに追加
	for i, fileID := range publishedFileIDs {
		params[fmt.Sprintf("publishedfileids[%d]", i)] = fileID
	}

	result, err := a.client.Request(ctx, "GET", EndpointPublishedFileDetails, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get published file details: %w", err)
	}

	return result, nil
}

// PublishedFileOptions はWorkshopファイル取得のオプション
type PublishedFileOptions struct {
	IncludeTags               bool `json:"includetags"`
	IncludeAdditionalPreviews bool `json:"includeadditionalpreviews"`
	IncludeChildren           bool `json:"includechildren"`
	IncludeKVTags             bool `json:"includekvtags"`
	IncludeVotes              bool `json:"includevotes"`
	ShortDescription          bool `json:"short_description"`
	IncludeForSaleData        bool `json:"includeforsaledata"`
	IncludeMetadata           bool `json:"includemetadata"`
	Language                  *int `json:"language"`
	ReturnPlaytimeStats       int  `json:"return_playtime_stats"`
	AppID                     *int `json:"appid"`
	StripDescriptionBBCode    bool `json:"strip_description_bbcode"`
	DesiredRevision           *int `json:"desired_revision"`
	IncludeReactions          bool `json:"includereactions"`
	AdminQuery                bool `json:"admin_query"`
}

// GetMultipleAppDetails は複数のアプリケーションの詳細情報を取得します
func (a *AppsService) GetMultipleAppDetails(ctx context.Context, appIDs []int, country, filters string) (map[int]*AppDetails, error) {
	if len(appIDs) == 0 {
		return make(map[int]*AppDetails), nil
	}

	if country == "" {
		country = "US"
	}

	if filters == "" {
		filters = "basic"
	}

	result := make(map[int]*AppDetails)

	// Steam Store APIは複数のApp IDを一度に処理できるが、
	// 安全のため一つずつ処理する
	for _, appID := range appIDs {
		if !validateAppID(appID) {
			continue
		}

		appDetails, err := a.GetAppDetails(ctx, appID, country, filters)
		if err != nil {
			// エラーが発生した場合はログに記録して続行
			fmt.Printf("Failed to get details for app %d: %v\n", appID, err)
			continue
		}

		result[appID] = appDetails
	}

	return result, nil
}

// GetAppList は利用可能なアプリケーション一覧を取得します
func (a *AppsService) GetAppList(ctx context.Context) (any, error) {
	result, err := a.client.Request(ctx, "GET", EndpointAppsAppList, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get app list: %w", err)
	}

	return result, nil
}

// GetServersAtAddress は指定されたアドレスのゲームサーバー情報を取得します
func (a *AppsService) GetServersAtAddress(ctx context.Context, addr string) (any, error) {
	if addr == "" {
		return nil, fmt.Errorf("address is required")
	}

	params := map[string]any{
		"addr": addr,
	}

	result, err := a.client.Request(ctx, "GET", EndpointAppsServersAtAddress, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get servers at address: %w", err)
	}

	return result, nil
}

// GetUpToDateCheck はアプリケーションのバージョンチェックを行います
func (a *AppsService) GetUpToDateCheck(ctx context.Context, appID int, version string) (any, error) {
	if !validateAppID(appID) {
		return nil, fmt.Errorf("invalid App ID: %d", appID)
	}

	if version == "" {
		return nil, fmt.Errorf("version is required")
	}

	params := map[string]any{
		"appid":   appID,
		"version": version,
	}

	result, err := a.client.Request(ctx, "GET", EndpointAppsUpToDateCheck, params)
	if err != nil {
		return nil, fmt.Errorf("failed to check up to date: %w", err)
	}

	return result, nil
}
