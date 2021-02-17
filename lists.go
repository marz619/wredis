package wredis

import (
	"github.com/garyburd/redigo/redis"
)

// LPush inserts the provided item(s) at the head of the list stored at key. For
// more information, see http://redis.io/commands/lpush.
func (w *impl) LPush(key string, items ...string) (int64, error) {
	if empty(key) {
		return int64Err("wredis: empty key")
	}
	if len(items) == 0 {
		return int64Err("must provide at least one item")
	}
	for _, i := range items {
		if empty(i) {
			return int64Err("an item cannot be empty")
		}
	}
	return w.Int64(func(conn redis.Conn) (int64, error) {
		args := redis.Args{}.Add(key).AddFlat(items)
		return redis.Int64(conn.Do("LPUSH", args...))
	})
}

// RPop removes and returns the last element of the list stored at key. For more
// information, see http://redis.io/commands/rpop.
func (w *impl) RPop(key string) (string, error) {
	if empty(key) {
		return stringErr("wredis: empty key")
	}
	return w.String(func(conn redis.Conn) (string, error) {
		args := redis.Args{}.Add(key)
		return redis.String(conn.Do("RPOP", args...))
	})
}

// LLen returns the length of the list stored at key. For more information, see
// http://redis.io/commands/llen.
func (w *impl) LLen(key string) (int64, error) {
	if empty(key) {
		return int64Err("wredis: empty key")
	}
	return w.Int64(func(conn redis.Conn) (int64, error) {
		args := redis.Args{}.Add(key)
		return redis.Int64(conn.Do("LLEN", args...))
	})
}
