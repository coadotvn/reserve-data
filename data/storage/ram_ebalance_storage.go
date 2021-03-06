package storage

import (
	"errors"
	"github.com/KyberNetwork/reserve-data/common"
	"sync"
)

type RamEBalanceStorage struct {
	mu      sync.RWMutex
	version int64
	data    map[int64]map[common.ExchangeID]common.EBalanceEntry
}

func NewRamEBalanceStorage() *RamEBalanceStorage {
	return &RamEBalanceStorage{
		mu:      sync.RWMutex{},
		version: 0,
		data:    map[int64]map[common.ExchangeID]common.EBalanceEntry{},
	}
}

func (self *RamEBalanceStorage) CurrentVersion() (int64, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.version, nil
}

func (self *RamEBalanceStorage) GetAllBalances(version int64) (map[common.ExchangeID]common.EBalanceEntry, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	all := self.data[version]
	if all == nil {
		return map[common.ExchangeID]common.EBalanceEntry{}, errors.New("Version doesn't exist")
	} else {
		return all, nil
	}
}

func (self *RamEBalanceStorage) StoreNewData(data map[common.ExchangeID]common.EBalanceEntry) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.version = self.version + 1
	self.data[self.version] = data
	delete(self.data, self.version-1)
	return nil
}
