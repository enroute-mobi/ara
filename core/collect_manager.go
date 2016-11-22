package core

type CollectManager struct {
	partners Partners
}

func NewCollectManager(partners Partners) *CollectManager {
	return &CollectManager{partners: partners}
}

func (manager *CollectManager) Partners() Partners {
	return manager.partners
}
