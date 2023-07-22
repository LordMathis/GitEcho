package storage

import "sync"

type StorageManager struct {
	mutex    sync.RWMutex
	storages map[string]Storage
}

func NewStorageManager() *StorageManager {
	return &StorageManager{
		mutex:    sync.RWMutex{},
		storages: make(map[string]Storage),
	}
}

func (sm *StorageManager) AddStorage(stor Storage) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.storages[stor.GetName()] = stor

}

func (sm *StorageManager) GetStorage(name string) Storage {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	return sm.storages[name]
}

func (sm *StorageManager) DeleteStorage(name string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.storages, name)
}
