package steam_api

// SteamError はSteam APIのエラーレスポンス
type SteamError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	StatusCode  int    `json:"-"`
}

func (e *SteamError) Error() string {
	return e.Description
}

// PlayerSummary はプレイヤーの基本情報
type PlayerSummary struct {
	SteamID                  string `json:"steamid"`
	CommunityVisibilityState int    `json:"communityvisibilitystate"`
	ProfileState             int    `json:"profilestate"`
	PersonaName              string `json:"personaname"`
	ProfileURL               string `json:"profileurl"`
	Avatar                   string `json:"avatar"`
	AvatarMedium             string `json:"avatarmedium"`
	AvatarFull               string `json:"avatarfull"`
	AvatarHash               string `json:"avatarhash"`
	LastLogoff               int64  `json:"lastlogoff"`
	PersonaState             int    `json:"personastate"`
	RealName                 string `json:"realname"`
	PrimaryClanID            string `json:"primaryclanid"`
	TimeCreated              int64  `json:"timecreated"`
	PersonaStateFlags        int    `json:"personastateflags"`
	LocCountryCode           string `json:"loccountrycode"`
	LocStateCode             string `json:"locstatecode"`
	LocCityID                int    `json:"loccityid"`
}

// PlayerSummariesResponse はGetPlayerSummariesのレスポンス
type PlayerSummariesResponse struct {
	Response struct {
		Players []PlayerSummary `json:"players"`
	} `json:"response"`
}

// ResolveVanityURLResponse はResolveVanityURLのレスポンス
type ResolveVanityURLResponse struct {
	Response struct {
		SteamID string `json:"steamid"`
		Success int    `json:"success"`
		Message string `json:"message"`
	} `json:"response"`
}

// Friend はフレンド情報
type Friend struct {
	SteamID      string `json:"steamid"`
	Relationship string `json:"relationship"`
	FriendSince  int64  `json:"friend_since"`
}

// FriendsListResponse はGetFriendListのレスポンス
type FriendsListResponse struct {
	FriendsList struct {
		Friends []Friend `json:"friends"`
	} `json:"friendslist"`
}

// OwnedGame は所有ゲーム情報
type OwnedGame struct {
	AppID                    int    `json:"appid"`
	Name                     string `json:"name"`
	PlaytimeForever          int    `json:"playtime_forever"`
	Playtime2Weeks           int    `json:"playtime_2weeks"`
	ImgIconURL               string `json:"img_icon_url"`
	ImgLogoURL               string `json:"img_logo_url"`
	PlaytimeWindowsForever   int    `json:"playtime_windows_forever"`
	PlaytimeMacForever       int    `json:"playtime_mac_forever"`
	PlaytimeLinuxForever     int    `json:"playtime_linux_forever"`
	PlaytimeDisconnected     int    `json:"playtime_disconnected"`
	RtimeLastPlayed          int64  `json:"rtime_last_played"`
	PlaytimeDeckForever      int    `json:"playtime_deck_forever"`
	HasCommunityVisibleStats bool   `json:"has_community_visible_stats"`
}

// OwnedGamesResponse はGetOwnedGamesのレスポンス
type OwnedGamesResponse struct {
	Response struct {
		GameCount int         `json:"game_count"`
		Games     []OwnedGame `json:"games"`
	} `json:"response"`
}

// RecentlyPlayedGame は最近プレイしたゲーム情報
type RecentlyPlayedGame struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	Playtime2Weeks  int    `json:"playtime_2weeks"`
	PlaytimeForever int    `json:"playtime_forever"`
	ImgIconURL      string `json:"img_icon_url"`
	ImgLogoURL      string `json:"img_logo_url"`
}

// RecentlyPlayedGamesResponse はGetRecentlyPlayedGamesのレスポンス
type RecentlyPlayedGamesResponse struct {
	Response struct {
		TotalCount int                  `json:"total_count"`
		Games      []RecentlyPlayedGame `json:"games"`
	} `json:"response"`
}

// SteamLevelResponse はGetSteamLevelのレスポンス
type SteamLevelResponse struct {
	Response struct {
		PlayerLevel int `json:"player_level"`
	} `json:"response"`
}

// Badge はバッジ情報
type Badge struct {
	BadgeID         int   `json:"badgeid"`
	Level           int   `json:"level"`
	CompletionTime  int64 `json:"completion_time"`
	XP              int   `json:"xp"`
	Scarcity        int   `json:"scarcity"`
	AppID           int   `json:"appid"`
	CommunityItemID int64 `json:"communityitemid"`
	BorderColor     int   `json:"border_color"`
}

// BadgesResponse はGetBadgesのレスポンス
type BadgesResponse struct {
	Response struct {
		Badges                     []Badge `json:"badges"`
		PlayerXP                   int     `json:"player_xp"`
		PlayerLevel                int     `json:"player_level"`
		PlayerXPNeededToLevelUp    int     `json:"player_xp_needed_to_level_up"`
		PlayerXPNeededCurrentLevel int     `json:"player_xp_needed_current_level"`
	} `json:"response"`
}

// PlayerBan はプレイヤーのBANステータス
type PlayerBan struct {
	SteamID          string `json:"SteamId"`
	CommunityBanned  bool   `json:"CommunityBanned"`
	VACBanned        bool   `json:"VACBanned"`
	NumberOfVACBans  int    `json:"NumberOfVACBans"`
	DaysSinceLastBan int    `json:"DaysSinceLastBan"`
	NumberOfGameBans int    `json:"NumberOfGameBans"`
	EconomyBan       string `json:"EconomyBan"`
}

// PlayerBansResponse はGetPlayerBansのレスポンス
type PlayerBansResponse struct {
	Players []PlayerBan `json:"players"`
}

// AppDetails はアプリケーション詳細情報
type AppDetails struct {
	Type                string                 `json:"type"`
	Name                string                 `json:"name"`
	SteamAppID          int                    `json:"steam_appid"`
	RequiredAge         int                    `json:"required_age"`
	IsFree              bool                   `json:"is_free"`
	DLC                 []int                  `json:"dlc"`
	DetailedDescription string                 `json:"detailed_description"`
	AboutTheGame        string                 `json:"about_the_game"`
	ShortDescription    string                 `json:"short_description"`
	SupportedLanguages  string                 `json:"supported_languages"`
	HeaderImage         string                 `json:"header_image"`
	Website             string                 `json:"website"`
	PCRequirements      map[string]interface{} `json:"pc_requirements"`
	MacRequirements     map[string]interface{} `json:"mac_requirements"`
	LinuxRequirements   map[string]interface{} `json:"linux_requirements"`
	Developers          []string               `json:"developers"`
	Publishers          []string               `json:"publishers"`
	PriceOverview       *PriceOverview         `json:"price_overview"`
	Packages            []int                  `json:"packages"`
	PackageGroups       []PackageGroup         `json:"package_groups"`
	Platforms           Platforms              `json:"platforms"`
	Metacritic          *Metacritic            `json:"metacritic"`
	Categories          []Category             `json:"categories"`
	Genres              []Genre                `json:"genres"`
	Screenshots         []Screenshot           `json:"screenshots"`
	Movies              []Movie                `json:"movies"`
	Recommendations     *Recommendations       `json:"recommendations"`
	Achievements        *Achievements          `json:"achievements"`
	ReleaseDate         ReleaseDate            `json:"release_date"`
	SupportInfo         SupportInfo            `json:"support_info"`
	Background          string                 `json:"background"`
	BackgroundRaw       string                 `json:"background_raw"`
	ContentDescriptors  ContentDescriptors     `json:"content_descriptors"`
}

// PriceOverview は価格情報
type PriceOverview struct {
	Currency                 string `json:"currency"`
	Initial                  int    `json:"initial"`
	Final                    int    `json:"final"`
	DiscountPercent          int    `json:"discount_percent"`
	InitialFormatted         string `json:"initial_formatted"`
	FinalFormatted           string `json:"final_formatted"`
	RecurringSub             int    `json:"recurring_sub"`
	RecurringSubDesc         string `json:"recurring_sub_desc"`
	IndividualPrice          int    `json:"individual_price"`
	IndividualPriceFormatted string `json:"individual_price_formatted"`
}

// PackageGroup はパッケージグループ情報
type PackageGroup struct {
	Name                    string `json:"name"`
	Title                   string `json:"title"`
	Description             string `json:"description"`
	SelectionText           string `json:"selection_text"`
	SaveText                string `json:"save_text"`
	DisplayType             int    `json:"display_type"`
	IsRecurringSubscription string `json:"is_recurring_subscription"`
}

// Platforms はプラットフォーム対応情報
type Platforms struct {
	Windows bool `json:"windows"`
	Mac     bool `json:"mac"`
	Linux   bool `json:"linux"`
}

// Metacritic はMetacriticスコア情報
type Metacritic struct {
	Score int    `json:"score"`
	URL   string `json:"url"`
}

// Category はカテゴリ情報
type Category struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

// Genre はジャンル情報
type Genre struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// Screenshot はスクリーンショット情報
type Screenshot struct {
	ID            int    `json:"id"`
	PathThumbnail string `json:"path_thumbnail"`
	PathFull      string `json:"path_full"`
}

// Movie は動画情報
type Movie struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	Webm      struct {
		Res480 string `json:"480"`
		Max    string `json:"max"`
	} `json:"webm"`
	Mp4 struct {
		Res480 string `json:"480"`
		Max    string `json:"max"`
	} `json:"mp4"`
	Highlight bool `json:"highlight"`
}

// Recommendations は推奨情報
type Recommendations struct {
	Total int `json:"total"`
}

// Achievements は実績情報
type Achievements struct {
	Total       int `json:"total"`
	Highlighted []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"highlighted"`
}

// ReleaseDate はリリース日情報
type ReleaseDate struct {
	ComingSoon bool   `json:"coming_soon"`
	Date       string `json:"date"`
}

// SupportInfo はサポート情報
type SupportInfo struct {
	URL   string `json:"url"`
	Email string `json:"email"`
}

// ContentDescriptors はコンテンツ記述子
type ContentDescriptors struct {
	IDs   []int  `json:"ids"`
	Notes string `json:"notes"`
}

// AppDetailsResponse はGetAppDetailsのレスポンス
type AppDetailsResponse map[string]struct {
	Success bool        `json:"success"`
	Data    *AppDetails `json:"data"`
}

// UserStats はユーザー統計情報
type UserStats struct {
	SteamID      string        `json:"steamID"`
	GameName     string        `json:"gameName"`
	Stats        []Stat        `json:"stats"`
	Achievements []Achievement `json:"achievements"`
}

// Stat は統計情報
type Stat struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// Achievement は実績情報
type Achievement struct {
	Name     string `json:"name"`
	Achieved int    `json:"achieved"`
}

// UserStatsResponse はGetUserStatsForGameのレスポンス
type UserStatsResponse struct {
	PlayerstatsResponse struct {
		SteamID      string        `json:"steamID"`
		GameName     string        `json:"gameName"`
		Stats        []Stat        `json:"stats"`
		Achievements []Achievement `json:"achievements"`
	} `json:"playerstats"`
}

// PlayerAchievements はプレイヤー実績情報
type PlayerAchievements struct {
	SteamID      string        `json:"steamID"`
	GameName     string        `json:"gameName"`
	Achievements []Achievement `json:"achievements"`
	Success      bool          `json:"success"`
}

// PlayerAchievementsResponse はGetPlayerAchievementsのレスポンス
type PlayerAchievementsResponse struct {
	PlayerstatsResponse PlayerAchievements `json:"playerstats"`
}
