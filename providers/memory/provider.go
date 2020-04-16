package memory

import (
	"fmt"
	"time"

	"github.com/fasthttp/session/v2"
	"github.com/savsgio/gotils"
)

var lastActiveTimeAttrKey = fmt.Sprintf("__session:lastActiveTime:%s__", gotils.RandBytes(make([]byte, 5)))

// New returns a new memory provider configured
func New(cfg Config) (*Provider, error) {
	p := &Provider{
		config: cfg,
		db:     new(session.Dict),
	}

	return p, nil
}

// Get sets the user session to the given store
func (p *Provider) Get(store *session.Store) error {
	data := p.db.GetBytes(store.GetSessionID())
	if data == nil {
		return nil
	}

	ptr := store.Ptr()
	ptr.D = append(ptr.D[:0], data.(session.Dict).D...)

	return nil
}

// Save saves the user session from the given store
func (p *Provider) Save(store *session.Store) error {
	store.Set(lastActiveTimeAttrKey, time.Now().Unix())
	p.db.SetBytes(store.GetSessionID(), store.GetAll())

	return nil
}

// Regenerate updates a user session with the new session id
// and sets the user session to the store
func (p *Provider) Regenerate(id []byte, newStore *session.Store) error {
	data := p.db.GetBytes(id)
	if data != nil {
		newStore.Ptr().D = data.(session.Dict).D
		p.Save(newStore)
		p.db.DelBytes(id)
	}

	return nil
}

// Destroy destroys the user session from the given id
func (p *Provider) Destroy(id []byte) error {
	p.db.DelBytes(id)

	return nil
}

// Count returns the total of users sessions stored
func (p *Provider) Count() int {
	return len(p.db.D)
}

// NeedGC indicates if the GC needs to be run
func (p *Provider) NeedGC() bool {
	return true
}

// GC destroys the expired user sessions
func (p *Provider) GC() {
	store := session.NewStore() // TODO: Get from pool
	ptr := store.Ptr()
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
