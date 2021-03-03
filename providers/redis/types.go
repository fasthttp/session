package redis

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/go-redis/redis/v8"
)

type CommonConfig struct {
	// Key prefix
	KeyPrefix string

	// Optional username.
	Username string

	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string

	// Maximum number of retries before giving up.
	// Default is to not retry failed commands.
	MaxRetries int

	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration

	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration

	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration

	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int

	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int

	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration

	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration

	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration

	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	// but idle connections are still discarded by the client
	// if IdleTimeout is set.
	IdleCheckFrequency time.Duration

	// TLS Config to use. When set TLS will be negotiated.
	TLSConfig *tls.Config
}

// Config provider settings
type Config struct {
	CommonConfig

	// Database to be selected after connecting to the server.
	DB int

	// The network type, either tcp or unix.
	// Default is tcp.
	Network string

	// host:port address.
	Addr string

	// Limiter interface used to implemented circuit breaker or rate limiter.
	Limiter redis.Limiter
}

// FailoverConfig provider settings
type FailoverConfig struct {
	CommonConfig

	// Database to be selected after connecting to the server.
	DB int

	// The sentinel master name.
	MasterName string

	// The sentinel nodes seed list (host:port).
	SentinelAddrs []string

	// The password for the sentinel connection if required (different to username/password).
	SentinelPassword string

	// Routes read-only commands to the closest node.
	RouteByLatency bool

	// Routes read-only commands in random order.
	RouteRandomly bool

	// Route read-only commands to slave nodes.
	SlaveOnly bool
}

// ClusterConfig provider settings
type ClusterConfig struct {
	CommonConfig

	Addrs []string

	// NewClient creates a cluster node client with provided name and options.
	NewClient func(opt *redis.Options) *redis.Client

	// The maximum number of retries before giving up. Command is retried
	// on network errors and MOVED/ASK redirects.
	// Default is 3 retries.
	MaxRedirects int

	// Routes read-only commands to the closest node.
	RouteByLatency bool

	// Routes read-only commands in random order.
	RouteRandomly bool

	// Optional function that returns cluster slots information.
	// It is useful to manually create cluster of standalone Redis servers
	// and load-balance read/write operations between master and slaves.
	// It can use service like ZooKeeper to maintain configuration information
	// and Cluster.ReloadState to manually trigger state reloading.
	ClusterSlots func(context.Context) ([]redis.ClusterSlot, error)
}

// Provider backend manager
type Provider struct {
	keyprefix string
	db        redis.Cmdable
}
