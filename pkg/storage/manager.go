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

func (sm *StorageManager) GetAllStorages() []Storage {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	storages := make([]Storage, 0, len(sm.storages))
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
