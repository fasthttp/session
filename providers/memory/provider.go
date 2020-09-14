package memory

import (
	"sync"
	"time"

	"github.com/fasthttp/session/v2"
	"github.com/savsgio/gotils"
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

// Get returns the data of the given session id
func (p *Provider) Get(id []byte) ([]byte, error) {
	val := p.db.GetBytes(id)
	if val == nil { // Not exist
		return nil, nil
	}

	item := val.(*item)

	return item.data, nil
}

// Save saves the session data and expiration from the given session id
func (p *Provider) Save(id, data []byte, expiration time.Duration) error {
	item := acquireItem()
	item.data = data
	item.lastActiveTime = time.Now().UnixNano()
	item.expiration = expiration

	p.db.SetBytes(id, item)

	return nil
}

// Regenerate updates the session id and expiration with the new session id
// of the the given current session id
func (p *Provider) Regenerate(id, newID []byte, expiration time.Duration) error {
	data := p.db.GetBytes(id)
	if data != nil {
		item := data.(*item)
		item.lastActiveTime = time.Now().UnixNano()
		item.expiration = expiration

		p.db.SetBytes(newID, item)
		p.db.DelBytes(id)
	}

	return nil
}

// Destroy destroys the session from the given id
func (p *Provider) Destroy(id []byte) error {
	val := p.db.GetBytes(id)
	if val == nil {
		return nil
	}

	p.db.DelBytes(id)
	releaseItem(val.(*item))

	return nil
}

// Count returns the total of stored sessions
func (p *Provider) Count() int {
	return len(p.db.D)
}

// NeedGC indicates if the GC needs to be run
func (p *Provider) NeedGC() bool {
	return true
}

// GC destroys the expired sessions
func (p *Provider) GC() {
	now := time.Now().UnixNano()

	for _, kv := range p.db.D {
		item := kv.Value.(*item)

		if item.expiration == 0 {
			continue
		}

		if now >= (item.lastActiveTime + item.expiration.Nanoseconds()) {
			p.Destroy(gotils.S2B(kv.Key))
		}
	}
}
