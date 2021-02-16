package wredis

import (
	"errors"
	"strings"

	"github.com/garyburd/redigo/redis"
)

// Del deletes one or more keys from Redis and returns a count of how many keys
// were actually deleted.
//
// See: http://redis.io/commands/del
func (w *poolClient) Del(keys ...string) (int64, error) {
	if len(keys) == 0 {
		return int64Err("wredis: no keys")
	}
	if any(keys, empty) {
		return int64Err("wredis: empty keys")
	}
	return w.Int64(func(conn redis.Conn) (int64, error) {
		args := redis.Args{}.AddFlat(keys)
		return redis.Int64(conn.Do("DEL", args...))
	})
}

// Delete is an alias for the Del method
func (w *poolClient) Delete(keys ...string) (int64, error) {
	return w.Del(keys...)
}

// DelPattern is a convenience method that Deletes *all* keys matching the
// provided pattern.
//
// See: http://redis.io/commands/keys
// See: http://redis.io/commands/del
func (w *poolClient) DelPattern(pattern string) (int64, error) {
	if !w.unsafe {
		return int64Err(unsafeErr("DelPattern").Error())
	}
	if strings.TrimSpace(pattern) == "" {
		return int64Err("empty pattern")
	}
	// fetch & delete keys
	keys, err := w.Keys(pattern)
	if err != nil {
		return int64Err(err.Error())
	}
	if len(keys) == 0 {
		return 0, nil
	}
	return w.Del(keys...)
}

// Exists checks for the existence of `key` in Redis. Note however, even though
// a variable number of keys can be passed to the EXISTS command since Redis
// 3.0.3, we will restrict this to a single key in order to return an absolute
// response regarding a key's existence.
//
// See: http://redis.io/commands/exists
func (w *poolClient) Exists(key string) (bool, error) {
	if empty(key) {
		return boolErr("wredis: empty key")
	}
	return w.Bool(func(conn redis.Conn) (bool, error) {
		return redis.Bool(conn.Do("EXISTS", key))
	})
}

// Expire sets a timeout of "seconds" on "key". If an error is encountered, it
// is returned. If the key doesn't exist or the timeout could not be set, then
// `false, nil` is returned. On success, `true, nil` is returned.
//
// See: http://redis.io/commands/expire
func (w *poolClient) Expire(key string, seconds int) (bool, error) {
	if empty(key) {
		return boolErr("wredis: empty key")
	}
	return w.Bool(func(conn redis.Conn) (bool, error) {
		return redis.Bool(conn.Do("EXPIRE", key, seconds))
	})
}

// Keys takes a pattern and returns any/all keys matching the pattern.
//
// See: http://redis.io/commands/keys
func (w *poolClient) Keys(pattern string) ([]string, error) {
	if empty(pattern) {
		return stringsErr("wredis: empty pattern")
	}
	return w.Strings(func(conn redis.Conn) ([]string, error) {
		return redis.Strings(conn.Do("KEYS", pattern))
	})
}

// Rename will rename some "from" to "to".
//
// See: `http://redis.io/commands/rename`
func (w *poolClient) Rename(from, to string) error {
	if empty(from) {
		return errors.New("wredis: empty from")
	}
	if empty(to) {
		return errors.New("wredis: empty to")
	}
	if from == to {
		return errors.New("wredis: from == wredis: empty key")
	}

	return w.ok("Rename", func(conn redis.Conn) (string, error) {
		return redis.String(conn.Do("RENAME", from, to))
	})
}
