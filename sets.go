package wredis

import "github.com/garyburd/redigo/redis"

// SAdd implements the SADD command. An error is returned if `members` is empty,
// otherwise it returns the number of members added to the Set at `dest`.
//
// See: http://redis.io/commands/sadd
func (w *poolClient) SAdd(dest string, members ...string) (int64, error) {
	if len(members) == 0 {
		return int64Err("wredis: no members")
	}
	return w.Int64(func(conn redis.Conn) (int64, error) {
		args := redis.Args{}.Add(dest).AddFlat(members)
		return redis.Int64(conn.Do("SADD", args...))
	})
}

// SCard returns the cardinality (size) of the Set at `key`.
//
// See: http://redis.io/commands/scard
func (w *poolClient) SCard(key string) (int64, error) {
	if empty(key) {
		return int64Err("wredis: empty key")
	}
	return w.Int64(func(conn redis.Conn) (int64, error) {
		return redis.Int64(conn.Do("SCARD", key))
	})
}

// SDiffStore executes the SDIFFSTORE command.
// See `http://redis.io/commands/sdiffstore`
func (w *poolClient) SDiffStore(dest string, keys ...string) (int64, error) {
	if empty(dest) {
		return int64Err("wredis: empty dest")
	}
	if len(keys) == 0 {
		return int64Err("wredis: no set keys")
	}
	if any(keys, empty) {
		return int64Err("wredis: empty set keys")
	}
	return w.Int64(func(conn redis.Conn) (int64, error) {
		return redis.Int64(conn.Do("SDIFFSTORE", redis.Args{}.Add(dest).AddFlat(keys)...))
	})
}

// SMembers returns the members of the set denoted by `key`.
// See: http://redis.io/commands/smembers
func (w *poolClient) SMembers(key string) ([]string, error) {
	if empty(key) {
		return stringsErr("wredis: empty key")
	}
	return w.Strings(func(conn redis.Conn) ([]string, error) {
		return redis.Strings(conn.Do("SMEMBERS", key))
	})
}

// SUnionStore implements the SUNIONSTORE command.
// See `http://redis.io/commands/sunionstore`
func (w *poolClient) SUnionStore(dest string, keys ...string) (int64, error) {
	if empty(dest) {
		return int64Err("wredis: empty dest")
	}
	if len(keys) == 0 {
		return int64Err("wredis: no set keys")
	}
	if any(keys, empty) {
		return int64Err("wredis: empty set keys")
	}
	return w.Int64(func(conn redis.Conn) (int64, error) {
		return redis.Int64(conn.Do("SUNIONSTORE", redis.Args{}.Add(dest).AddFlat(keys)...))
	})
}
