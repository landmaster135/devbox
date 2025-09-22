package steam_api

const (
	// APIBaseURL はSteam Web APIのベースURL
	APIBaseURL = "https://api.steampowered.com"

	// APIAppDetailsURL はSteam Store APIのアプリ詳細取得URL
	APIAppDetailsURL = "https://store.steampowered.com/api/appdetails"

	// APIAppSearchURL はSteam Store APIのアプリ検索URL
	APIAppSearchURL = "https://store.steampowered.com/search/suggest"

	// APIAppWishlistURL はWishlist取得用のURL
	APIAppWishlistURL = "https://api.steampowered.com/IWishlistService/GetWishlist/v1/"

	// User関連のエンドポイント
	EndpointUserPlayerSummaries  = "/ISteamUser/GetPlayerSummaries/v2/"
	EndpointUserFriendList       = "/ISteamUser/GetFriendList/v1/"
	EndpointUserPlayerBans       = "/ISteamUser/GetPlayerBans/v1/"
	EndpointUserResolveVanityURL = "/ISteamUser/ResolveVanityURL/v1/"

	// Player Service関連のエンドポイント
	EndpointPlayerRecentlyPlayedGames    = "/IPlayerService/GetRecentlyPlayedGames/v1/"
	EndpointPlayerOwnedGames             = "/IPlayerService/GetOwnedGames/v1/"
	EndpointPlayerSteamLevel             = "/IPlayerService/GetSteamLevel/v1/"
	EndpointPlayerBadges                 = "/IPlayerService/GetBadges/v1/"
	EndpointPlayerCommunityBadgeProgress = "/IPlayerService/GetCommunityBadgeProgress/v1/"

	// User Stats関連のエンドポイント
	EndpointUserStatsForGame       = "/ISteamUserStats/GetUserStatsForGame/v2/"
	EndpointUserPlayerAchievements = "/ISteamUserStats/GetPlayerAchievements/v1/"

	// Apps関連のエンドポイント
	EndpointAppsAppList          = "/ISteamApps/GetAppList/v2/"
	EndpointAppsServersAtAddress = "/ISteamApps/GetServersAtAddress/v1/"
	EndpointAppsUpToDateCheck    = "/ISteamApps/UpToDateCheck/v1/"

	// Published File Service関連のエンドポイント
	EndpointPublishedFileDetails = "/IPublishedFileService/GetDetails/v1/"

	// Game Servers Service関連のエンドポイント
	EndpointGameServersAccountPublicInfo = "/IGameServersService/GetAccountPublicInfo/v1/"

	HTTPStatusCode400              = 400
	HTTPStatusCode500              = 500
	ThresholdOfIDCountsForSteamAPI = 100
	DigitsOfSteamID                = 17
)
