package domain

import "time"

// AniListResponse はAniList GraphQL APIのレスポンス構造体
type AniListResponse struct {
	Data   *MediaListCollectionData `json:"data"`
	Errors []GraphQLError           `json:"errors,omitempty"`
}

// GraphQLError はGraphQLエラーの構造体
type GraphQLError struct {
	Message   string                 `json:"message"`
	Locations []GraphQLErrorLocation `json:"locations,omitempty"`
	Path      []interface{}          `json:"path,omitempty"`
}

// GraphQLErrorLocation はGraphQLエラーの位置情報
type GraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// MediaListCollectionData はMediaListCollectionクエリのデータ部分
type MediaListCollectionData struct {
	MediaListCollection *MediaListCollection `json:"MediaListCollection"`
}

// MediaListCollection はメディアリストコレクション
type MediaListCollection struct {
	Lists []MediaList `json:"lists"`
}

// MediaList はメディアリスト
type MediaList struct {
	Entries []MediaListEntry `json:"entries"`
}

// MediaListEntry はメディアリストエントリ
type MediaListEntry struct {
	Media       *Media      `json:"media"`
	Score       int         `json:"score"`
	Status      string      `json:"status"`
	Progress    int         `json:"progress"`
	CompletedAt *FuzzyDate  `json:"completedAt"`
	Notes       string      `json:"notes"`
	UpdatedAt   int64       `json:"updatedAt"`
}

// Media はメディア情報
type Media struct {
	ID         int      `json:"id"`
	Title      *Title   `json:"title"`
	CoverImage *Image   `json:"coverImage"`
	SiteURL    string   `json:"siteUrl"`
	Studios    *Studios `json:"studios"`
}

// Title はタイトル情報
type Title struct {
	Native string `json:"native"`
}

// Image は画像情報
type Image struct {
	ExtraLarge string `json:"extraLarge"`
}

// Studios はスタジオ情報
type Studios struct {
	Nodes []Studio `json:"nodes"`
}

// Studio はスタジオ
type Studio struct {
	Name string `json:"name"`
}

// FuzzyDate は曖昧な日付
type FuzzyDate struct {
	Year  *int `json:"year"`
	Month *int `json:"month"`
	Day   *int `json:"day"`
}

// AnimeInfo は整形されたアニメ情報
type AnimeInfo struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Score         int       `json:"score"`
	Status        string    `json:"status"`
	Progress      int       `json:"progress"`
	CompletedAt   string    `json:"completed_at"`
	Notes         string    `json:"notes"`
	CoverImageURL string    `json:"cover_image_url"`
	SiteURL       string    `json:"site_url"`
	Studio        string    `json:"studio"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GraphQLRequest はGraphQLリクエストの構造体
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// QueryAnimeRequest はquery-anime操作のリクエスト
type QueryAnimeRequest struct {
	Username string
	UserID   *int
}
