package memcache

// Save save store
func (mcs *Store) Save() error {
	data := mcs.GetAll()
	value, err := provider.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	item := acquireItem()
	item.Key = provider.getMemCacheSessionKey(mcs.GetSessionID())
	item.Value = value
	item.Expiration = provider.expiration

	err = provider.db.Set(item)

	releaseItem(item)

	return err
}
