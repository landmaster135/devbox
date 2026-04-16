package domain

type Artifact struct {
	ConID     string        `json:"con_id"`
	PageTitle string        `json:"page_title"`
	Category  string        `json:"category"`
	OutputURL string        `json:"output_url"`
	Tags      []ArtifactTag `json:"tags"`
}

type ArtifactTag struct {
	PageTitle string `json:"page_title"`
}
