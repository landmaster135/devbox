package valkey

import (
	"sync"
)

var (
	instance *Server
	once     sync.Once
	mutex    sync.Mutex
)

// GetServer はシングルトンパターンでValkeyサーバーのインスタンスを取得します
func GetServer(address string) (*Server, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var err error
	once.Do(func() {
		instance, err = NewServer(address)
		if err == nil && instance != nil {
			err = instance.Start()
		}
	})

	return instance, err
}

// ShutdownServer はValkeyサーバーを安全に停止します
func ShutdownServer() {
	mutex.Lock()
	defer mutex.Unlock()

	if instance != nil {
		instance.Stop()
		instance = nil
	}
}
