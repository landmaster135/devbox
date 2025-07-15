package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          Utility Function Tests                            ##
// #==============================================================#
func TestBraveSearchService_getOrDefault_Normal(t *testing.T) {
	service := NewBraveSearchService()

	// 通常の値
	result := service.getOrDefault("test", "default")
	assert.Equal(t, "test", result)

	// 空の値
	result = service.getOrDefault("", "default")
	assert.Equal(t, "default", result)
}

// #==============================================================#
// ##          Format Local Results Tests                        ##
// #==============================================================#
func TestBraveSearchService_formatLocalResults_Normal(t *testing.T) {
	service := NewBraveSearchService()

	// テストデータの作成
	poisData := BravePoiResponse{
		Results: []BraveLocation{
			{
				ID:   "location1",
				Name: "Test Location 1",
				Address: struct {
					StreetAddress   string `json:"streetAddress,omitempty"`
					AddressLocality string `json:"addressLocality,omitempty"`
					AddressRegion   string `json:"addressRegion,omitempty"`
					PostalCode      string `json:"postalCode,omitempty"`
				}{
					StreetAddress:   "123 Test St",
					AddressLocality: "Test City",
					AddressRegion:   "Test Region",
					PostalCode:      "12345",
				},
				Phone: "123-456-7890",
				Rating: struct {
					RatingValue float64 `json:"ratingValue,omitempty"`
					RatingCount int     `json:"ratingCount,omitempty"`
				}{
					RatingValue: 4.5,
					RatingCount: 100,
				},
				OpeningHours: []string{"Mon-Fri: 9AM-5PM", "Sat-Sun: Closed"},
				PriceRange:   "$$",
			},
		},
	}

	descData := BraveDescription{
		Descriptions: map[string]string{
			"location1": "This is a test location 1",
		},
	}

	// テストの実行
	result := service.formatLocalResults(poisData, descData)

	// 結果の検証
	assert.Contains(t, result, "Name: Test Location 1")
	assert.Contains(t, result, "Address: 123 Test St, Test City, Test Region, 12345")
	assert.Contains(t, result, "Phone: 123-456-7890")
	assert.Contains(t, result, "Rating: 4.5 (100 reviews)")
	assert.Contains(t, result, "Price Range: $$")
	assert.Contains(t, result, "Hours: Mon-Fri: 9AM-5PM, Sat-Sun: Closed")
	assert.Contains(t, result, "Description: This is a test location 1")
}

func TestBraveSearchService_formatLocalResults_EmptyResults(t *testing.T) {
	service := NewBraveSearchService()

	// 空のデータ
	emptyPoisData := BravePoiResponse{
		Results: []BraveLocation{},
	}
	descData := BraveDescription{}

	result := service.formatLocalResults(emptyPoisData, descData)
	assert.Equal(t, "No local results found", result)
}

func TestBraveSearchService_formatLocalResults_MissingData(t *testing.T) {
	service := NewBraveSearchService()

	// データが不完全な場合
	poisData := BravePoiResponse{
		Results: []BraveLocation{
			{
				ID:   "location1",
				Name: "Test Location 1",
				// 他のフィールドは空
			},
		},
	}

	descData := BraveDescription{
		Descriptions: map[string]string{},
	}

	result := service.formatLocalResults(poisData, descData)

	// デフォルト値が使用されることを確認
	assert.Contains(t, result, "Name: Test Location 1")
	assert.Contains(t, result, "Address: N/A")
	assert.Contains(t, result, "Phone: N/A")
	assert.Contains(t, result, "Rating: N/A")
	assert.Contains(t, result, "Price Range: N/A")
	assert.Contains(t, result, "Hours: N/A")
	assert.Contains(t, result, "Description: No description available")
}

func TestBraveSearchService_formatLocalResults_MultipleLocations(t *testing.T) {
	service := NewBraveSearchService()

	// 複数のロケーションデータ
	poisData := BravePoiResponse{
		Results: []BraveLocation{
			{
				ID:   "location1",
				Name: "Location 1",
				Address: struct {
					StreetAddress   string `json:"streetAddress,omitempty"`
					AddressLocality string `json:"addressLocality,omitempty"`
					AddressRegion   string `json:"addressRegion,omitempty"`
					PostalCode      string `json:"postalCode,omitempty"`
				}{
					StreetAddress: "123 First St",
				},
				Phone: "111-111-1111",
			},
			{
				ID:   "location2",
				Name: "Location 2",
				Address: struct {
					StreetAddress   string `json:"streetAddress,omitempty"`
					AddressLocality string `json:"addressLocality,omitempty"`
					AddressRegion   string `json:"addressRegion,omitempty"`
					PostalCode      string `json:"postalCode,omitempty"`
				}{
					AddressLocality: "Second City",
				},
				Rating: struct {
					RatingValue float64 `json:"ratingValue,omitempty"`
					RatingCount int     `json:"ratingCount,omitempty"`
				}{
					RatingValue: 3.8,
					RatingCount: 50,
				},
			},
		},
	}

	descData := BraveDescription{
		Descriptions: map[string]string{
			"location1": "First location description",
		},
	}

	// テストの実行
	result := service.formatLocalResults(poisData, descData)

	// 結果の検証
	assert.Contains(t, result, "Name: Location 1")
	assert.Contains(t, result, "Address: 123 First St")
	assert.Contains(t, result, "Phone: 111-111-1111")
	assert.Contains(t, result, "Description: First location description")
	assert.Contains(t, result, "---") // セパレータ
	assert.Contains(t, result, "Name: Location 2")
	assert.Contains(t, result, "Address: Second City")
	assert.Contains(t, result, "Rating: 3.8 (50 reviews)")
	assert.Contains(t, result, "Description: No description available")
}

func TestBraveSearchService_formatLocalResults_PartialAddress(t *testing.T) {
	service := NewBraveSearchService()

	// 部分的な住所データ
	poisData := BravePoiResponse{
		Results: []BraveLocation{
			{
				ID:   "location1",
				Name: "Test Location",
				Address: struct {
					StreetAddress   string `json:"streetAddress,omitempty"`
					AddressLocality string `json:"addressLocality,omitempty"`
					AddressRegion   string `json:"addressRegion,omitempty"`
					PostalCode      string `json:"postalCode,omitempty"`
				}{
					AddressLocality: "Test City",
					PostalCode:      "12345",
					// StreetAddressとAddressRegionは空
				},
			},
		},
	}

	descData := BraveDescription{}

	// テストの実行
	result := service.formatLocalResults(poisData, descData)

	// 結果の検証
	assert.Contains(t, result, "Address: Test City, 12345")
	assert.NotContains(t, result, "Address: , Test City, , 12345") // 空の部分は含まれない
}

func TestBraveSearchService_formatLocalResults_ZeroRating(t *testing.T) {
	service := NewBraveSearchService()

	// 評価が0の場合
	poisData := BravePoiResponse{
		Results: []BraveLocation{
			{
				ID:   "location1",
				Name: "Test Location",
				Rating: struct {
					RatingValue float64 `json:"ratingValue,omitempty"`
					RatingCount int     `json:"ratingCount,omitempty"`
				}{
					RatingValue: 0.0,
					RatingCount: 0,
				},
			},
		},
	}

	descData := BraveDescription{}

	// テストの実行
	result := service.formatLocalResults(poisData, descData)

	// 結果の検証
	assert.Contains(t, result, "Rating: N/A")
}
