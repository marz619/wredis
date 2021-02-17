package wredis

import (
	"errors"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Append the value to the string denoted by key. If the key does not exist,
// then it is created and set as an empty string.
//
// See: https://redis.io/commands/append
func (w *impl) Append(key, value string) (int64, error) {
	if empty(key) {
		return int64Err("wredis: empty key")
	}
	return w.Int64(func(conn redis.Conn) (int64, error) {
		return redis.Int64(conn.Do("APPEND", key, value))
	})
}

// Appends is a convenience wrapper that accepts a variadic number of strings
// and pre-joins them before calling Append.
//
// See: https://redis.io/commands/append
func (w *impl) Appends(key, sep string, values ...string) (int64, error) {
	return w.Append(key, strings.Join(values, sep))
}

// Get retrieves the string value for some key.
//
// See: http://redis.io/commands/get
func (w *impl) Get(key string) (string, error) {
	if empty(key) {
		return stringErr("wredis: empty key")
	}
	return w.String(func(conn redis.Conn) (string, error) {
		return redis.String(conn.Do("GET", key))
	})
}

// MGet returns the values of all provided keys. For a key that does not exist,
// an empty string is returned.
//
// See: http://redis.io/commands/mget.
func (w *impl) MGet(keys ...string) ([]string, error) {
	if any(keys, empty) {
		return stringsErr("wredis: empty keys")
	}
	return w.Strings(func(conn redis.Conn) ([]string, error) {
		args := redis.Args{}.AddFlat(keys)
		return redis.Strings(conn.Do("MGET", args...))
	})
}

// Incr increments the number stored at key by one.
//
// See: http://redis.io/commands/incr
func (w *impl) Incr(key string) (int64, error) {
	if empty(key) {
		return int64Err("wredis: empty key")
	}
	return w.Int64(func(conn redis.Conn) (int64, error) {
		return redis.Int64(conn.Do("INCR", key))
	})
}

// Set a string value to some key.
//
// See: http://redis.io/commands/set
func (w *impl) Set(key, value string) error {
	if empty(key) {
		return errors.New("wredis: empty key")
	}
	return w.ok("Set", func(conn redis.Conn) (string, error) {
		return redis.String(conn.Do("SET", key, value))
	})
}

// SetEx sets key's to value with an expiry time measured in seconds.
//
// See: http://redis.io/commands/setex
func (w *impl) SetEx(key, value string, seconds uint) error {
	if empty(key) {
		return errors.New("wredis: empty key")
	}
	if seconds == uint(0) {
		return errors.New("wredis: one second expiry")
	}
	return w.ok("SetEx", func(conn redis.Conn) (string, error) {
		args := redis.Args{}.Add(key).Add(seconds).Add(value)
		return redis.String(conn.Do("SETEX", args...))
	})
}

// SetExDuration is a convenience method that calls SetEx, but sets the expiry
// value using a time.Duration.
//
// See: http://redis.io/commands/setex
func (w *impl) SetExDuration(k, v string, d time.Duration) error {
	return w.SetEx(k, v, uint(d.Seconds()))
}
