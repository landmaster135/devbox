package steam_api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// #==============================================================#
// ##       Mocks for UsersService                               ##
// #==============================================================#
// MockUsersService はUsersServiceのモック実装
type MockUsersService struct {
	SearchUserFunc                 func(ctx context.Context, vanityURL string) (*PlayerSummary, error)
	GetUserDetailsFunc             func(ctx context.Context, steamID string) (*PlayerSummary, error)
	GetMultipleUserDetailsFunc     func(ctx context.Context, steamIDs []string) ([]PlayerSummary, error)
	GetUserFriendsListFunc         func(ctx context.Context, steamID string, enriched bool) ([]Friend, error)
	GetUserRecentlyPlayedGamesFunc func(ctx context.Context, steamID string) ([]RecentlyPlayedGame, error)
	GetOwnedGamesFunc              func(ctx context.Context, steamID string, includeAppInfo, includeFreeGames bool) ([]OwnedGame, error)
	GetUserSteamLevelFunc          func(ctx context.Context, steamID string) (int, error)
	GetUserBadgesFunc              func(ctx context.Context, steamID string) ([]Badge, error)
	GetPlayerBansFunc              func(ctx context.Context, steamID string) (*PlayerBan, error)
	GetSteamIDFunc                 func(ctx context.Context, vanityURL string) (string, error)
	GetCommunityBadgeProgressFunc  func(ctx context.Context, steamID string, badgeID int) (any, error)
	GetAccountPublicInfoFunc       func(ctx context.Context, steamID string) (any, error)
	GetClientFunc                  func() ClientInterface
}

// SearchUser はモックのSearchUserメソッド
func (m *MockUsersService) SearchUser(ctx context.Context, vanityURL string) (*PlayerSummary, error) {
	if m.SearchUserFunc != nil {
		return m.SearchUserFunc(ctx, vanityURL)
	}
	return nil, nil
}

// GetUserDetails はモックのGetUserDetailsメソッド
func (m *MockUsersService) GetUserDetails(ctx context.Context, steamID string) (*PlayerSummary, error) {
	if m.GetUserDetailsFunc != nil {
		return m.GetUserDetailsFunc(ctx, steamID)
	}
	return nil, nil
}

// GetMultipleUserDetails はモックのGetMultipleUserDetailsメソッド
func (m *MockUsersService) GetMultipleUserDetails(ctx context.Context, steamIDs []string) ([]PlayerSummary, error) {
	if m.GetMultipleUserDetailsFunc != nil {
		return m.GetMultipleUserDetailsFunc(ctx, steamIDs)
	}
	return nil, nil
}

// GetUserFriendsList はモックのGetUserFriendsListメソッド
func (m *MockUsersService) GetUserFriendsList(ctx context.Context, steamID string, enriched bool) ([]Friend, error) {
	if m.GetUserFriendsListFunc != nil {
		return m.GetUserFriendsListFunc(ctx, steamID, enriched)
	}
	return nil, nil
}

// GetUserRecentlyPlayedGames はモックのGetUserRecentlyPlayedGamesメソッド
func (m *MockUsersService) GetUserRecentlyPlayedGames(ctx context.Context, steamID string) ([]RecentlyPlayedGame, error) {
	if m.GetUserRecentlyPlayedGamesFunc != nil {
		return m.GetUserRecentlyPlayedGamesFunc(ctx, steamID)
	}
	return nil, nil
}

// GetOwnedGames はモックのGetOwnedGamesメソッド
func (m *MockUsersService) GetOwnedGames(ctx context.Context, steamID string, includeAppInfo, includeFreeGames bool) ([]OwnedGame, error) {
	if m.GetOwnedGamesFunc != nil {
		return m.GetOwnedGamesFunc(ctx, steamID, includeAppInfo, includeFreeGames)
	}
	return nil, nil
}

// GetUserSteamLevel はモックのGetUserSteamLevelメソッド
func (m *MockUsersService) GetUserSteamLevel(ctx context.Context, steamID string) (int, error) {
	if m.GetUserSteamLevelFunc != nil {
		return m.GetUserSteamLevelFunc(ctx, steamID)
	}
	return 0, nil
}

// GetUserBadges はモックのGetUserBadgesメソッド
func (m *MockUsersService) GetUserBadges(ctx context.Context, steamID string) ([]Badge, error) {
	if m.GetUserBadgesFunc != nil {
		return m.GetUserBadgesFunc(ctx, steamID)
	}
	return nil, nil
}

// GetPlayerBans はモックのGetPlayerBansメソッド
func (m *MockUsersService) GetPlayerBans(ctx context.Context, steamID string) (*PlayerBan, error) {
	if m.GetPlayerBansFunc != nil {
		return m.GetPlayerBansFunc(ctx, steamID)
	}
	return nil, nil
}

// GetSteamID はモックのGetSteamIDメソッド
func (m *MockUsersService) GetSteamID(ctx context.Context, vanityURL string) (string, error) {
	if m.GetSteamIDFunc != nil {
		return m.GetSteamIDFunc(ctx, vanityURL)
	}
	return "", nil
}

// GetCommunityBadgeProgress はモックのGetCommunityBadgeProgressメソッド
func (m *MockUsersService) GetCommunityBadgeProgress(ctx context.Context, steamID string, badgeID int) (any, error) {
	if m.GetCommunityBadgeProgressFunc != nil {
		return m.GetCommunityBadgeProgressFunc(ctx, steamID, badgeID)
	}
	return nil, nil
}

// GetAccountPublicInfo はモックのGetAccountPublicInfoメソッド
func (m *MockUsersService) GetAccountPublicInfo(ctx context.Context, steamID string) (any, error) {
	if m.GetAccountPublicInfoFunc != nil {
		return m.GetAccountPublicInfoFunc(ctx, steamID)
	}
	return nil, nil
}

// GetClient はモックのGetClientメソッド
func (m *MockUsersService) GetClient() ClientInterface {
	if m.GetClientFunc != nil {
		return m.GetClientFunc()
	}
	return nil
}

// #==============================================================#
// ##       Interfaces for UsersService                          ##
// #==============================================================#
// UsersServiceInterface はUsers関連のAPIを抽象化
type UsersServiceInterface interface {
	SearchUser(ctx context.Context, vanityURL string) (*PlayerSummary, error)
	GetUserDetails(ctx context.Context, steamID string) (*PlayerSummary, error)
	GetMultipleUserDetails(ctx context.Context, steamIDs []string) ([]PlayerSummary, error)
	GetUserFriendsList(ctx context.Context, steamID string, enriched bool) ([]Friend, error)
	GetUserRecentlyPlayedGames(ctx context.Context, steamID string) ([]RecentlyPlayedGame, error)
	GetOwnedGames(ctx context.Context, steamID string, includeAppInfo, includeFreeGames bool) ([]OwnedGame, error)
	GetUserSteamLevel(ctx context.Context, steamID string) (int, error)
	GetUserBadges(ctx context.Context, steamID string) ([]Badge, error)
	GetPlayerBans(ctx context.Context, steamID string) (*PlayerBan, error)
	GetSteamID(ctx context.Context, vanityURL string) (string, error)
	GetCommunityBadgeProgress(ctx context.Context, steamID string, badgeID int) (any, error)
	GetAccountPublicInfo(ctx context.Context, steamID string) (any, error)
	GetClient() ClientInterface
}

// #==============================================================#
// ##       Implementations for UsersService                     ##
// #==============================================================#
// UsersService はSteam Users APIのサービス
type UsersService struct {
	client ClientInterface
}

// NewUsersService は新しいUsersServiceを作成します
func NewUsersService(client ClientInterface) *UsersService {
	return &UsersService{
		client: client,
	}
}

// GetClient はクライアントを返します
func (u *UsersService) GetClient() ClientInterface {
	return u.client
}

// SearchUser はユーザー名でユーザーを検索します
func (u *UsersService) SearchUser(ctx context.Context, vanityURL string) (*PlayerSummary, error) {
	if vanityURL == "" {
		return nil, fmt.Errorf("vanity URL is required")
	}

	// まずVanity URLをSteam IDに変換
	steamID, err := u.GetSteamID(ctx, vanityURL)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve vanity URL: %w", err)
	}

	// Steam IDでユーザー詳細を取得
	return u.GetUserDetails(ctx, steamID)
}

// GetUserDetails はSteam IDでユーザー詳細を取得します
func (u *UsersService) GetUserDetails(ctx context.Context, steamID string) (*PlayerSummary, error) {
	if !validateSteamID(steamID) {
		return nil, fmt.Errorf("invalid Steam ID: %s", steamID)
	}

	params := map[string]any{
		"steamids": steamID,
	}

	result, err := u.GetClient().Request(ctx, "GET", EndpointUserPlayerSummaries, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	// レスポンスをパース
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var response PlayerSummariesResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Response.Players) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &response.Response.Players[0], nil
}

// GetMultipleUserDetails は複数のSteam IDでユーザー詳細を取得します
func (u *UsersService) GetMultipleUserDetails(ctx context.Context, steamIDs []string) ([]PlayerSummary, error) {
	if len(steamIDs) == 0 {
		return []PlayerSummary{}, nil
	}

	var allPlayers []PlayerSummary

	// Steam APIは一度に100個までのIDしか処理できないため、チャンクに分割
	chunks := chunkSteamIDs(steamIDs, ThresholdOfIDCountsForSteamAPI)

	for _, chunk := range chunks {
		validIDs := make([]string, 0, len(chunk))
		for _, id := range chunk {
			if validateSteamID(id) {
				validIDs = append(validIDs, id)
			}
		}

		if len(validIDs) == 0 {
			continue
		}

		params := map[string]any{
			"steamids": strings.Join(validIDs, ","),
		}

		result, err := u.GetClient().Request(ctx, "GET", EndpointUserPlayerSummaries, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get user details: %w", err)
		}

		// レスポンスをパース
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}

		var response PlayerSummariesResponse
		if err := json.Unmarshal(jsonBytes, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		allPlayers = append(allPlayers, response.Response.Players...)
	}

	return allPlayers, nil
}

// GetUserFriendsList はユーザーのフレンドリストを取得します
func (u *UsersService) GetUserFriendsList(ctx context.Context, steamID string, enriched bool) ([]Friend, error) {
	if !validateSteamID(steamID) {
		return nil, fmt.Errorf("invalid Steam ID: %s", steamID)
	}

	params := map[string]any{
		"steamid": steamID,
	}

	result, err := u.GetClient().Request(ctx, "GET", EndpointUserFriendList, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get friends list: %w", err)
	}

	// レスポンスをパース
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var response FriendsListResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	friends := response.FriendsList.Friends

	// enrichedが有効な場合、フレンドの詳細情報も取得
	if enriched && len(friends) > 0 {
		friendIDs := make([]string, len(friends))
		for i, friend := range friends {
			friendIDs[i] = friend.SteamID
		}

		playerDetails, err := u.GetMultipleUserDetails(ctx, friendIDs)
		if err != nil {
			return friends, nil // エラーが発生してもフレンドリストは返す
		}

		// フレンド情報と詳細情報をマージ
		playerMap := make(map[string]PlayerSummary)
		for _, player := range playerDetails {
			playerMap[player.SteamID] = player
		}

		// ここでは簡単にFriendスライスを返すが、実際にはEnrichedFriend構造体を作成することも可能
	}

	return friends, nil
}

// GetUserRecentlyPlayedGames は最近プレイしたゲームを取得します
func (u *UsersService) GetUserRecentlyPlayedGames(ctx context.Context, steamID string) ([]RecentlyPlayedGame, error) {
	if !validateSteamID(steamID) {
		return nil, fmt.Errorf("invalid Steam ID: %s", steamID)
	}

	params := map[string]any{
		"steamid": steamID,
	}

	result, err := u.GetClient().Request(ctx, "GET", EndpointPlayerRecentlyPlayedGames, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get recently played games: %w", err)
	}

	// レスポンスをパース
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var response RecentlyPlayedGamesResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Response.Games, nil
}

// GetOwnedGames は所有ゲーム一覧を取得します
func (u *UsersService) GetOwnedGames(ctx context.Context, steamID string, includeAppInfo, includeFreeGames bool) ([]OwnedGame, error) {
	if !validateSteamID(steamID) {
		return nil, fmt.Errorf("invalid Steam ID: %s", steamID)
	}

	params := map[string]any{
		"steamid":                   steamID,
		"include_appinfo":           includeAppInfo,
		"include_played_free_games": includeFreeGames,
	}

	result, err := u.GetClient().Request(ctx, "GET", EndpointPlayerOwnedGames, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned games: %w", err)
	}

	// レスポンスをパース
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var response OwnedGamesResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Response.Games, nil
}

// GetUserSteamLevel はユーザーのSteamレベルを取得します
func (u *UsersService) GetUserSteamLevel(ctx context.Context, steamID string) (int, error) {
	if !validateSteamID(steamID) {
		return 0, fmt.Errorf("invalid Steam ID: %s", steamID)
	}

	params := map[string]any{
		"steamid": steamID,
	}

	result, err := u.GetClient().Request(ctx, "GET", EndpointPlayerSteamLevel, params)
	if err != nil {
		return 0, fmt.Errorf("failed to get steam level: %w", err)
	}

	// レスポンスをパース
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal response: %w", err)
	}

	var response SteamLevelResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Response.PlayerLevel, nil
}

// GetUserBadges はユーザーのバッジ情報を取得します
func (u *UsersService) GetUserBadges(ctx context.Context, steamID string) ([]Badge, error) {
	if !validateSteamID(steamID) {
		return nil, fmt.Errorf("invalid Steam ID: %s", steamID)
	}

	params := map[string]any{
		"steamid": steamID,
	}

	result, err := u.GetClient().Request(ctx, "GET", EndpointPlayerBadges, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get badges: %w", err)
	}

	// レスポンスをパース
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var response BadgesResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Response.Badges, nil
}

// GetPlayerBans はプレイヤーのBAN情報を取得します
func (u *UsersService) GetPlayerBans(ctx context.Context, steamID string) (*PlayerBan, error) {
	if !validateSteamID(steamID) {
		return nil, fmt.Errorf("invalid Steam ID: %s", steamID)
	}

	params := map[string]any{
		"steamids": steamID,
	}

	result, err := u.GetClient().Request(ctx, "GET", EndpointUserPlayerBans, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get player bans: %w", err)
	}

	// レスポンスをパース
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var response PlayerBansResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Players) == 0 {
		return nil, fmt.Errorf("player ban info not found")
	}

	return &response.Players[0], nil
}

// GetSteamID はVanity URLからSteam IDを取得します
func (u *UsersService) GetSteamID(ctx context.Context, vanityURL string) (string, error) {
	if vanityURL == "" {
		return "", fmt.Errorf("vanity URL is required")
	}

	params := map[string]any{
		"vanityurl": vanityURL,
	}

	result, err := u.GetClient().Request(ctx, "GET", EndpointUserResolveVanityURL, params)
	if err != nil {
		return "", fmt.Errorf("failed to resolve vanity URL: %w", err)
	}

	// レスポンスをパース
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	var response ResolveVanityURLResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Response.Success != 1 {
		// より詳細なエラーメッセージを提供
		switch response.Response.Success {
		case 42:
			return "", fmt.Errorf("vanity URL '%s' does not exist or is not set.  Note: Make sure not the display name", vanityURL)
		default:
			message := response.Response.Message
			if message == "" {
				message = "Unknown error"
			}
			return "", fmt.Errorf("failed to resolve vanity URL '%s': %s (error code: %d)", vanityURL, message, response.Response.Success)
		}
	}

	return response.Response.SteamID, nil
}

// GetCommunityBadgeProgress はコミュニティバッジの進行状況を取得します
func (u *UsersService) GetCommunityBadgeProgress(ctx context.Context, steamID string, badgeID int) (any, error) {
	if !validateSteamID(steamID) {
		return nil, fmt.Errorf("invalid Steam ID: %s", steamID)
	}

	if badgeID <= 0 {
		return nil, fmt.Errorf("invalid badge ID: %d", badgeID)
	}

	params := map[string]any{
		"steamid": steamID,
		"badgeid": badgeID,
	}

	result, err := u.GetClient().Request(ctx, "GET", EndpointPlayerCommunityBadgeProgress, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get community badge progress: %w", err)
	}

	return result, nil
}

// GetAccountPublicInfo はアカウントの公開情報を取得します
func (u *UsersService) GetAccountPublicInfo(ctx context.Context, steamID string) (any, error) {
	if !validateSteamID(steamID) {
		return nil, fmt.Errorf("invalid Steam ID: %s", steamID)
	}

	params := map[string]any{
		"steamid": steamID,
	}

	result, err := u.GetClient().Request(ctx, "GET", EndpointGameServersAccountPublicInfo, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get account public info: %w", err)
	}

	return result, nil
}
