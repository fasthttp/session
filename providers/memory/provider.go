package memory

import (
	"sync"
	"time"

	"github.com/fasthttp/session/v2"
)

var itemPool = &sync.Pool{
	New: func() interface{} {
		return new(item)
	},
}

func acquireItem() *item {
	return itemPool.Get().(*item)
}

func releaseItem(item *item) {
	item.data = item.data[:0]
	item.lastActiveTime = 0
	item.expiration = 0

	itemPool.Put(item)
}

// New returns a new memory provider configured
func New(cfg Config) (*Provider, error) {
	p := &Provider{
		config: cfg,
		db:     new(session.Dict),
	}

	return p, nil
}

// Get sets the user session to the given store
func (p *Provider) Get(id []byte) ([]byte, error) {
	val := p.db.GetBytes(id)
	if val == nil { // Not exist
		return nil, nil
	}

	item := val.(*item)

	return item.data, nil
}

// Save saves the user session from the given store
func (p *Provider) Save(id, data []byte, expiration time.Duration) error {
	item := acquireItem()
	item.data = data
	item.lastActiveTime = time.Now().Unix()
	item.expiration = expiration

	p.db.SetBytes(id, item)

	return nil
}

// Regenerate updates a user session with the new session id
// and sets the user session to the store
func (p *Provider) Regenerate(id, newID []byte, expiration time.Duration) error {
	data := p.db.GetBytes(id)
	if data != nil {
		item := data.(*item)
		item.lastActiveTime = time.Now().Unix()
		item.expiration = expiration

		p.db.SetBytes(newID, item)
		p.db.DelBytes(id)
	}

	return nil
}

// Destroy destroys the user session from the given id
func (p *Provider) Destroy(id []byte) error {
	val := p.db.GetBytes(id)
	if val == nil {
		return nil
	}

	p.db.DelBytes(id)
	releaseItem(val.(*item))

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
	now := time.Now().Unix()

	for _, kv := range p.db.D {
		item := kv.Value.(*item)

		if item.expiration == 0 {
			continue
		}

		if now >= (item.lastActiveTime + int64(item.expiration.Seconds())) {
			p.Destroy(kv.Key)
		}
	}
}
