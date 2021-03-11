module github.com/fasthttp/session/v2

go 1.12

require (
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/go-redis/redis/v8 v8.3.4 // Don't upgrade to keep go1.12 compatibility
	github.com/go-sql-driver/mysql v1.5.0
	github.com/lib/pq v1.10.0
	github.com/mattn/go-sqlite3 v1.14.6
	github.com/savsgio/dictpool v0.0.0-20210217113430-85d3b37fb239
	github.com/savsgio/gotils v0.0.0-20210225112730-595c7e5a8a7a
	github.com/valyala/bytebufferpool v1.0.0
	github.com/valyala/fasthttp v1.22.0
)
