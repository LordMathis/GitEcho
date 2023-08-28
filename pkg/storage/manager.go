package storage

import "sync"

type StorageManager struct {
	mutex    sync.RWMutex
	storages map[string]BaseStorage
}

func NewStorageManager() *StorageManager {
	return &StorageManager{
		mutex:    sync.RWMutex{},
		storages: make(map[string]BaseStorage),
	}
}

func (sm *StorageManager) AddStorage(stor BaseStorage) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.storages[stor.Name] = stor

}

func (sm *StorageManager) GetStorage(name string) BaseStorage {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	return sm.storages[name]
}

func (sm *StorageManager) GetAllStorages() []BaseStorage {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	storages := make([]BaseStorage, 0, len(sm.storages))
	for _, stor := range sm.storages {
		storages = append(storages, stor)
	}

	return storages
}

func (sm *StorageManager) DeleteStorage(name string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.storages, name)
}
