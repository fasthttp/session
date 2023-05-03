//go:build go1.19
// +build go1.19

package redis

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config provider settings
type Config struct {
	// Key prefix
	KeyPrefix string

	// Pointer to the logger interface.
	Logger Logger

	// The network type, either tcp or unix.
	// Default is tcp.
	Network string

	// host:port address.
	Addr string

	// Optional username.
	Username string

	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string

	// Database to be selected after connecting to the server.
	DB int

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

	// Maximum number of idle connections.
	MaxIdleConns int

	// Deprecated: This field has been renamed to ConnMaxLifetime
	MaxConnAge time.Duration

	// Maximum amount of time a connection may be reused.
	// Expired connections may be closed lazily before reuse.
	// If <= 0, connections are not closed due to a connection's age.
	// Default is to not close idle connections.
	ConnMaxLifetime time.Duration

	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration

	// Deprecated: This field has been renamed to ConnMaxIdleTime
	IdleTimeout time.Duration

	// Deprecated: This field has been removed in favor of MaxIdleConns
	IdleCheckFrequency time.Duration

	// Maximum amount of time a connection may be idle.
	// Should be less than server's timeout.
	// Expired connections may be closed lazily before reuse.
	// If d <= 0, connections are not closed due to a connection's idle time.
	// Default is 30 minutes. -1 disables idle timeout check.
	ConnMaxIdleTime time.Duration

	// TLS Config to use. When set TLS will be negotiated.
	TLSConfig *tls.Config

	// Limiter interface used to implemented circuit breaker or rate limiter.
	Limiter redis.Limiter
}

// FailoverConfig provider settings.
type FailoverConfig struct {
	// Key prefix
	KeyPrefix string

	// Pointer to the logger interface.
	Logger Logger

	// Optional username.
	Username string

	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string

	// Database to be selected after connecting to the server.
	DB int

	// The sentinel master name.
	MasterName string

	// The sentinel nodes seed list (host:port).
	SentinelAddrs []string

	// The username to use for the sentinel connection if required. If specified, the Redis
	// client will attempt to authenticate via ACL authentication. If not specified, the
	// client will use requirepass-style authentication.
	SentinelUsername string

	// The password for the sentinel connection if required (different to username/password).
	SentinelPassword string

	// Routes read-only commands to the closest node. Only relevant with NewFailoverCluster.
	RouteByLatency bool

	// Routes read-only commands in random order. Only relevant with NewFailoverCluster.
	RouteRandomly bool

	// Deprecated: his field has been renamed to ReplicaOnly
	SlaveOnly bool

	// Route all commands to replica read-only nodes.
	ReplicaOnly bool

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

	// Maximum number of idle connections.
	MaxIdleConns int

	// Deprecated: This field has been renamed to ConnMaxLifetime
	MaxConnAge time.Duration

	// Maximum amount of time a connection may be reused.
	// Expired connections may be closed lazily before reuse.
	// If <= 0, connections are not closed due to a connection's age.
	// Default is to not close idle connections.
	ConnMaxLifetime time.Duration

	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration

	// Deprecated: This field has been renamed to ConnMaxIdleTime
	IdleTimeout time.Duration

	// Maximum amount of time a connection may be idle.
	// Should be less than server's timeout.
	// Expired connections may be closed lazily before reuse.
	// If d <= 0, connections are not closed due to a connection's idle time.
	// Default is 30 minutes. -1 disables idle timeout check.
	ConnMaxIdleTime time.Duration

	// TLS Config to use. When set TLS will be negotiated.
	TLSConfig *tls.Config
}

// Provider backend manager
type Provider struct {
	keyPrefix string
	db        redis.Cmdable
}

// Logger implements the upstream redis internal Logger interface.
type Logger interface {
	Printf(ctx context.Context, format string, v ...interface{})
}
