package domain

// EmbeddingConfig は埋め込みサービスへの接続情報を表す。
type EmbeddingConfig struct {
	Host           string
	Port           int
	Model          string
	TimeoutSeconds int
}

// CreateCollectionParams は create-collection 操作の入力を表す。
type CreateCollectionParams struct {
	DBHost         string
	DBPort         int
	CollectionName string
	Size           int
}

// DescribeCollectionParams は describe-collection 操作の入力を表す。
type DescribeCollectionParams struct {
	DBHost         string
	DBPort         int
	CollectionName string
}

// DeleteCollectionParams は delete-collection 操作の入力を表す。
type DeleteCollectionParams struct {
	DBHost         string
	DBPort         int
	CollectionName string
}

// ListCollectionsParams は list-collections 操作の入力を表す。
type ListCollectionsParams struct {
	DBHost string
	DBPort int
}

// InputPayloadPair は upsert-texts 用に 1 件の input と payload をまとめたもの。
type InputPayloadPair struct {
	Input   string
	Payload map[string]string
}

// UpsertTextsParams は upsert-texts 操作の入力。
type UpsertTextsParams struct {
	DBHost         string
	DBPort         int
	CollectionName string
	Embedding      EmbeddingConfig
	Pairs          []InputPayloadPair
}

// PayloadFilter は payload の単一条件を表す。
type PayloadFilter struct {
	Key   string
	Value string
}

// QueryPointsParams は query-points 操作の入力。
type QueryPointsParams struct {
	DBHost         string
	DBPort         int
	CollectionName string
	Embedding      EmbeddingConfig
	Input          string
	PayloadFilters []PayloadFilter
	Limit          int
}

// OverwritePayloadParams は overwrite-payload 操作の入力。
type OverwritePayloadParams struct {
	DBHost         string
	DBPort         int
	CollectionName string
	Payload        map[string]string
	PointIDs       []string
	Filters        OverwriteFilters
}

// OverwriteFilters は overwrite-payload 用の複数条件を表す。
type OverwriteFilters struct {
	Must           []PayloadFilter
	MustNot        []PayloadFilter
	Should         []PayloadFilter
	MinShouldCount int
}
