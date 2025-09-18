package steam_api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// UsersService はSteam Users APIのサービス
type UsersService struct {
	client *Client
}

// NewUsersService は新しいUsersServiceを作成します
func NewUsersService(client *Client) *UsersService {
	return &UsersService{
		client: client,
	}
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

	result, err := u.client.Request(ctx, "GET", "/ISteamUser/GetPlayerSummaries/v2/", params)
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
	chunks := chunkSteamIDs(steamIDs, 100)

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

		result, err := u.client.Request(ctx, "GET", "/ISteamUser/GetPlayerSummaries/v2/", params)
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

	result, err := u.client.Request(ctx, "GET", "/ISteamUser/GetFriendList/v1/", params)
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

	result, err := u.client.Request(ctx, "GET", "/IPlayerService/GetRecentlyPlayedGames/v1/", params)
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

	result, err := u.client.Request(ctx, "GET", "/IPlayerService/GetOwnedGames/v1/", params)
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

	result, err := u.client.Request(ctx, "GET", "/IPlayerService/GetSteamLevel/v1/", params)
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

	result, err := u.client.Request(ctx, "GET", "/IPlayerService/GetBadges/v1/", params)
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

	result, err := u.client.Request(ctx, "GET", "/ISteamUser/GetPlayerBans/v1/", params)
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

	result, err := u.client.Request(ctx, "GET", "/ISteamUser/ResolveVanityURL/v1/", params)
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

	result, err := u.client.Request(ctx, "GET", "/IPlayerService/GetCommunityBadgeProgress/v1/", params)
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

	result, err := u.client.Request(ctx, "GET", "/IGameServersService/GetAccountPublicInfo/v1/", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get account public info: %w", err)
	}

	return result, nil
}
