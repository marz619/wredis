package wredis

import (
	"errors"

	"github.com/garyburd/redigo/redis"
)

// Echo echoes the message.
//
// See: https://redis.io/commands/echo
func (w *poolClient) Echo(msg string) (string, error) {
	return w.match("Echo", msg, func(conn redis.Conn) (string, error) {
		return redis.String(conn.Do("ECHO", redis.Args{}.Add(msg)...))
	})
}

// Ping returns PONG if no message is provided or the message.
//
// See: https://redis.io/commands/ping
func (w *poolClient) Ping(msg ...string) (string, error) {
	if len(msg) > 1 {
		return stringErr("wredis: ping single message")
	}

	exp := "PONG"
	if len(msg) == 1 {
		exp = msg[0]
	}

	return w.match("Ping", exp, func(conn redis.Conn) (string, error) {
		return redis.String(conn.Do("PING", redis.Args{}.AddFlat(msg)...))
	})
}

// Quit asks the server to close the connection.
//
// See: https://redis.io/commands/quit
func (w *poolClient) Quit() error {
	return w.ok("Quit", func(conn redis.Conn) (string, error) {
		return redis.String(conn.Do("QUIT"))
	})
}

// Select selects the Database specified by the parameter. We use an unsigned
// int because Redis databases are numbered using a zero (0) based index.
//
// See: https://redis.io/commands/select
//
// NOTE: Implementation Details
//
// 1. Select only modifes the "current" Connection
// 2. We will return a new Wredis object that is *NOT* Selectable
func (w *poolClient) Select(db uint) (Wredis, error) {
	// Cannot call select in Cluster mode
	if w.cfg.unselectable {
		return nil, errors.New("wredis: no select")
	}

	// this will return a Wredis whose underlying pool contains 1 allowable
	// active connection that will never idle or timeout
	cfg, err := w.cfg.Copy(
		DB(db),
		IdleTimeout(0),
		MaxActive(1),
		MaxConnLifetime(0),
		MaxIdle(0),
		Wait(false),
		unselectable(),
	)
	if err != nil {
		return nil, err
	}

	// create and return our single connection Pool
	pool, err := newPoolClient(cfg)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
