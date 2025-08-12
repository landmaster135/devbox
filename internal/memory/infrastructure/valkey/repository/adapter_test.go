package valkey

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DataStoreAdapterTestSuite struct {
	mockDataStore *MockDataStore
	adapter       DataRepository
	ctx           context.Context
}

func setupDataStoreAdapterTest(t *testing.T) *DataStoreAdapterTestSuite {
	// 実際のDataStoreAdapterのテストは、実際のDataStoreインスタンスが必要なため、
	// ここではモックを使用したテスト用のアダプターを作成します
	mockDataStore := &MockDataStore{}

	// テスト用のアダプターを直接作成
	adapter := &DataStoreAdapter{
		dataStore: nil, // 実際のテストではモックを使用
		addr:      "localhost:6379",
	}

	return &DataStoreAdapterTestSuite{
		mockDataStore: mockDataStore,
		adapter:       adapter,
		ctx:           context.Background(),
	}
}

func TestNewDataStoreAdapter_Normal(t *testing.T) {
	// 実際のDataStoreを使用したテストは外部依存があるため、
	// ここでは基本的な構造体の作成をテストします
	addr := "localhost:6379"

	// 実際のテストでは外部依存があるため、構造体の初期化のみテスト
	adapter := &DataStoreAdapter{
		addr: addr,
	}

	assert.NotNil(t, adapter)
	assert.Equal(t, addr, adapter.GetServerAddress())
}

func TestNewDataStoreAdapter_NilDataStore(t *testing.T) {
	adapter, err := NewDataStoreAdapter(nil, "localhost:6379")

	assert.Error(t, err)
	assert.Nil(t, adapter)
	assert.Contains(t, err.Error(), "dataStore is nil")
}

func TestNewDataStoreAdapter_EmptyAddr(t *testing.T) {
	mockDataStore := &MockDataStore{}

	adapter, err := NewDataStoreAdapter(mockDataStore, "")

	assert.Error(t, err)
	assert.Nil(t, adapter)
	assert.Contains(t, err.Error(), "addr is empty")
}

func TestNewDataStoreAdapter_InvalidDataStoreType(t *testing.T) {
	invalidDataStore := "not a datastore"

	adapter, err := NewDataStoreAdapter(invalidDataStore, "localhost:6379")

	assert.Error(t, err)
	assert.Nil(t, adapter)
	assert.Contains(t, err.Error(), "dataStore is not *infraValkey.DataStore")
}

func TestDataStoreAdapter_StartServer_Normal(t *testing.T) {
	suite := setupDataStoreAdapterTest(t)

	err := suite.adapter.StartServer()

	assert.NoError(t, err) // アダプターでは何もしないのでエラーなし
}

func TestDataStoreAdapter_StopServer_Normal(t *testing.T) {
	suite := setupDataStoreAdapterTest(t)

	err := suite.adapter.StopServer()

	assert.NoError(t, err) // アダプターでは何もしないのでエラーなし
}

func TestDataStoreAdapter_IsServerRunning_Normal(t *testing.T) {
	suite := setupDataStoreAdapterTest(t)

	isRunning := suite.adapter.IsServerRunning()

	assert.True(t, isRunning) // アダプターでは常にtrueを返す
}

func TestDataStoreAdapter_GetServerAddress_Normal(t *testing.T) {
	suite := setupDataStoreAdapterTest(t)
	expectedAddr := "localhost:6379"

	addr := suite.adapter.GetServerAddress()

	assert.Equal(t, expectedAddr, addr)
}

// 実際のDataStoreAdapterの実装をテストするためのテストケース
// これらのテストは実際のdataStoreがnilのため、パニックが発生することを確認します

func TestDataStoreAdapter_RealImplementation_GetKeys_Panic(t *testing.T) {
	adapter := &DataStoreAdapter{
		dataStore: nil,
		addr:      "localhost:6379",
	}

	// dataStoreがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		adapter.GetKeys(context.Background(), "test:*")
	})
}

func TestDataStoreAdapter_RealImplementation_GetValue_Panic(t *testing.T) {
	adapter := &DataStoreAdapter{
		dataStore: nil,
		addr:      "localhost:6379",
	}

	// dataStoreがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		adapter.GetValue(context.Background(), "test_key")
	})
}

func TestDataStoreAdapter_RealImplementation_GetValueAsByte_Panic(t *testing.T) {
	adapter := &DataStoreAdapter{
		dataStore: nil,
		addr:      "localhost:6379",
	}

	// dataStoreがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		adapter.GetValueAsByte(context.Background(), "test_key")
	})
}

func TestDataStoreAdapter_RealImplementation_GetType_Panic(t *testing.T) {
	adapter := &DataStoreAdapter{
		dataStore: nil,
		addr:      "localhost:6379",
	}

	// dataStoreがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		adapter.GetType(context.Background(), "test_key")
	})
}

func TestDataStoreAdapter_RealImplementation_SetValue_Panic(t *testing.T) {
	adapter := &DataStoreAdapter{
		dataStore: nil,
		addr:      "localhost:6379",
	}

	// dataStoreがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		adapter.SetValue(context.Background(), "test_key", []byte("test_value"))
	})
}

func TestDataStoreAdapter_RealImplementation_DeleteKey_Panic(t *testing.T) {
	adapter := &DataStoreAdapter{
		dataStore: nil,
		addr:      "localhost:6379",
	}

	// dataStoreがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		adapter.DeleteKey(context.Background(), "test_key")
	})
}
