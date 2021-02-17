package wredis

import (
	"github.com/garyburd/redigo/redis"
)

// DBSize returns the number of keys in the selected database
//
// See: https://redis.io/commands/dbsize
func (w *impl) DBSize() (int64, error) {
	return w.Int64(func(conn redis.Conn) (int64, error) {
		return redis.Int64(conn.Do("DBSIZE"))
	})
}

// FlushAll deletes all the keys from all the db's on the Redis Server. This is
// very dangerous method to use in a production; do so at your own peril.
//
// See: http://redis.io/commands/flushall
func (w *impl) FlushAll() error {
	if !w.unsafe {
		return unsafeErr("FlushAll")
	}
	// all hands to battle stations!
	return w.ok("FlushAll", func(conn redis.Conn) (string, error) {
		return redis.String(conn.Do("FLUSHALL"))
	})
}

// FlushDB deletes all the keys from the configured DB
//
// See: http://redis.io/commands/flushdb
func (w *impl) FlushDB() error {
	if !w.unsafe {
		return unsafeErr("FlushDB")
	}
	// all stop, red alert!
	return w.ok("FlushDB", func(conn redis.Conn) (string, error) {
		return redis.String(conn.Do("FlUSHDB"))
	})
}
