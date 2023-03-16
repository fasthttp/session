package session

import (
	"testing"
	"time"
)

func Test_NewStore(t *testing.T) {
	store := NewStore()

	if store == nil {
		t.Error("Store is nil")
	}
}

func TestStore_GetSetDelete(t *testing.T) {
	store := NewStore()

	key := "key"
	value := "value"

	store.Set(key, value)

	if v := store.Get(key); v == nil {
		t.Errorf("Store.Get() == %v, want %v", v, value)
	}

	if v := store.Get("fake"); v != nil {
		t.Errorf("Store.Get() == %v, want %v", v, nil)
	}

	store.Delete(key)

	if v := store.Get(key); v != nil {
		t.Errorf("Store.Get() (after delete) == %v, want %v", v, nil)
	}
}

func TestStore_GetSetDeleteBytes(t *testing.T) {
	store := NewStore()

	key := []byte("key")
	value := "value"

	store.SetBytes(key, value)

	if v := store.GetBytes(key); v == nil {
		t.Errorf("Store.GetBytes() == %v, want %v", v, value)
	}

	if v := store.GetBytes([]byte("fake")); v != nil {
		t.Errorf("Store.GetBytes() == %v, want %v", v, nil)
	}

	store.DeleteBytes(key)

	if v := store.GetBytes(key); v != nil {
		t.Errorf("Store.GetBytes() (after delete) == %v, want %v", v, nil)
	}
}

func TestStore_Ptr(t *testing.T) {
	store := NewStore()

	if store.Ptr() != &store.data {
		t.Errorf("Store.Ptr() ==  %p, want %p", store.Ptr(), store.data)
	}
}

func TestStore_Flush(t *testing.T) {
	store := NewStore()
	store.Set("k1", "v1")
	store.Set("k2", "v2")
	store.Set("k3", "v3")

	store.Flush()

	if len(store.data.KV) > 0 {
		t.Error("Store is not flushed")
	}
}

func TestStore_SetGetSessionID(t *testing.T) {
	store := NewStore()

	id := []byte("1234abcd567trfg")

	store.SetSessionID(id)

	if v := store.GetSessionID(); string(v) != string(id) {
		t.Errorf("Store.GetSessionID() == %s, want %s", v, id)
	}
}

func TestStore_SetGetHasExpiration(t *testing.T) {
	store := NewStore()

	expiration := 10 * time.Second

	if err := store.SetExpiration(expiration); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if v := store.GetExpiration(); v != expiration {
		t.Errorf("Store.GetExpiration() == %d, want %d", v, expiration)
	}

	if !store.HasExpirationChanged() {
		t.Errorf("Store.HasExpirationChanged() == %v, want %v", false, true)
	}
}

func TestStore_Reset(t *testing.T) {
	store := NewStore()
	store.defaultExpiration = 10
	store.SetSessionID([]byte("af123443z"))
	store.Set("k", "v")

	store.Reset()

	if len(store.data.KV) > 0 || len(store.sessionID) > 0 || store.defaultExpiration != 0 {
		t.Error("Store is not reseted")
	}
}
