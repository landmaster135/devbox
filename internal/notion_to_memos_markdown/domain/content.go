package domain

type Content struct {
	ConID        string       `json:"con_id"`
	Category     string       `json:"category"`
	PageTitle    string       `json:"page_title"`
	OwningStatus string       `json:"owning_status"`
	Color        string       `json:"color"`
	BoughtAt     string       `json:"bought_at"`
	Score        int          `json:"score"`
	Price        int          `json:"price"`
	URL          string       `json:"url"`
	Tags         []ContentTag `json:"tags"`
}

type ContentTag struct {
	PageTitle string `json:"page_title"`
}
