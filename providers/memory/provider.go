package memory

import (
	"fmt"
	"time"

	"github.com/fasthttp/session"
	"github.com/savsgio/gotils"
)

var lastActiveTimeAttrKey = fmt.Sprintf("__session:lastActiveTime:%s__", gotils.RandBytes(make([]byte, 5)))

// NewProvider new memory provider
func New(cfg Config) (*Provider, error) {
	p := &Provider{
		config: cfg,
		db:     new(session.Dict),
	}

	return p, nil
}

// Get get session store by id
func (p *Provider) Get(store *session.Store) error {
	data := p.db.GetBytes(store.GetSessionID())
	if data == nil {
		return nil
	}

	ptr := store.DataPointer()
	ptr.D = append(ptr.D[:0], data.(session.Dict).D...)

	return nil
}

func (p *Provider) Save(store *session.Store) error {
	store.Set(lastActiveTimeAttrKey, time.Now().Unix())
	p.db.SetBytes(store.GetSessionID(), store.GetAll())

	return nil
}

// Regenerate regenerate session
func (p *Provider) Regenerate(id []byte, newStore *session.Store) error {
	data := p.db.GetBytes(id)
	if data != nil {
		newStore.DataPointer().D = data.(session.Dict).D
		p.Save(newStore)
		p.db.DelBytes(id)
	}

	return nil
}

// Destroy destroy session by sessionID
func (p *Provider) Destroy(id []byte) error {
	p.db.DelBytes(id)

	return nil
}

// Count session values count
func (p *Provider) Count() int {
	return len(p.db.D)
}

// NeedGC need gc
func (p *Provider) NeedGC() bool {
	return true
}

// GC session garbage collection
func (p *Provider) GC() {
	store := session.NewStore() // TODO: Get from pool
	ptr := store.DataPointer()
	now := time.Now().Unix()

	for _, kv := range p.db.D {
		ptr.D = append(ptr.D[:0], kv.Value.(session.Dict).D...)

		expiration := store.GetExpiration()
		// Do not expire the session if expiration is set to 0
		if expiration == 0 {
			continue
		}

		lastActiveTime := store.Get(lastActiveTimeAttrKey).(int64)

		if now >= (lastActiveTime + int64(expiration)) {
			p.Destroy(kv.Key)
		}
	}
}
