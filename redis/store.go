package redis

// Save save store
func (rs *Store) Save() error {
	data := rs.GetAll()
	b, err := provider.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	err = provider.db.Set(provider.getRedisSessionKey(rs.GetSessionID()), b, provider.expiration).Err()

	return err
}
