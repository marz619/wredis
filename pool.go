package wredis

import (
	"fmt"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

// interface check
var _ Wredis = &poolClient{}

// poolClient is a simple wrapper around the redis.Pool, which implements the
// Wredis interface
//
// See: http://redis.io/commands
type poolClient struct {
	cfg    Config      // Config this was intialised with
	pool   *redis.Pool // the underlying redis connection pool
	unsafe bool        // safe poolClient?

	mu     sync.RWMutex
	counts map[string]int // command counts
}

// get the command counts
func (w *poolClient) stats() CMDCounts {
	w.mu.RLock()
	defer w.mu.Unlock()

	counts := make(CMDCounts, len(w.counts))
	for k, v := range w.counts {
		counts[k] = v
	}
	return counts
}

// increment provided command count
func (w *poolClient) inc(cmd string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.counts[cmd]++
}

// Close will close the *redis.Pool
func (w *poolClient) Close() error {
	return w.pool.Close()
}

// Conn returns a redis.Conn from the underlying pool
func (w *poolClient) Conn() (redis.Conn, error) {
	// get a connection from the pool
	conn := w.pool.Get()
	// check the connection was established without error
	if err := conn.Err(); err != nil {
		return nil, err
	}
	return conn, nil
}

// Stats contains poolClient statistics.
type Stats struct {
	Stats  redis.PoolStats
	Counts CMDCounts
}

// CMDCounts is a simple wrapper around map[string]int
type CMDCounts map[string]int

// Count the count from some command, -1 if the command is not found.
func (cc CMDCounts) Count(cmd string) int {
	if v, ok := cc[cmd]; ok {
		return v
	}
	return -1
}

// Stats returns the current statstics.
func (w *poolClient) Stats() Stats {
	w.mu.Lock()
	defer w.mu.Unlock()

	return Stats{
		Stats:  w.pool.Stats(),
		Counts: w.stats(),
	}
}

var nilErr error = nil

// new returns a "safe" *poolClient poolClient with the configured options
func newPoolClient(cfg Config) (*poolClient, error) {
	// set up the *redis.Pool with our Config
	pool := &redis.Pool{
		MaxActive:       cfg.MaxActive,
		MaxConnLifetime: time.Duration(cfg.MaxConnLifetime),
		MaxIdle:         cfg.MaxIdle,
		IdleTimeout:     time.Duration(cfg.IdleTimeout),
		Dial:            cfg.Dialer(cfg),
		TestOnBorrow:    cfg.TestOnBorrower(cfg),
		Wait:            cfg.Wait,
	}

	return &poolClient{
		cfg:    cfg,
		pool:   pool,
		counts: make(map[string]int),
	}, nil
}

// Safe returns a "safe" *poolClient poolClient configured with the provided options.
func Safe(opts ...Option) (Wredis, error) {
	cfg, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}
	return newPoolClient(cfg)
}

// Unsafe returns an "unsafe" *poolClient poolClient configured with the provided
// options. The "unsafe"ness allows usage of certain methods that could be
// harmful if accidentally invoked in a production environment (e.g. FlushAll)
func Unsafe(opts ...Option) (Wredis, error) {
	cfg, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}
	w, err := newPoolClient(cfg)
	if err != nil {
		return nil, err
	}
	w.unsafe = true
	return w, nil
}

func (w *poolClient) ok(cmd string, f stringFunc) error {
	_, err := w.match(cmd, "OK", f)
	return err
}

const matchErrFmt = `wredis: %s expected "%s" response, got: "%s"`

func (w *poolClient) match(cmd, m string, f stringFunc) (string, error) {
	res, err := w.String(f)
	if err != nil {
		return stringErr(err.Error())
	}

	if res != m {
		return stringErr(fmt.Sprintf(matchErrFmt, cmd, m, res))
	}

	return res, nil
}

type (
	boolFunc    func(redis.Conn) (bool, error)
	int64Func   func(redis.Conn) (int64, error)
	stringFunc  func(redis.Conn) (string, error)
	stringsFunc func(redis.Conn) ([]string, error)
)

// Close is a default connection closer
func Close(conn redis.Conn) error {
	return conn.Close()
}

// Bool is a helper function to execute any series of commands
// on a redis.Conn that returns a bool response
func (w *poolClient) Bool(f boolFunc) (bool, error) {
	conn, err := w.Conn()
	if err != nil {
		return boolErr(err.Error())
	}
	defer Close(conn)
	return f(conn)
}

// Int64 is a helper function to execute any series of commands
// on a redis.Conn that return an int64 response
func (w *poolClient) Int64(f int64Func) (int64, error) {
	conn, err := w.Conn()
	if err != nil {
		return int64Err(err.Error())
	}
	defer Close(conn)
	return f(conn)
}

// String is a helper function to execute any series of commands
// on a redis.Conn that return a string response
func (w *poolClient) String(f stringFunc) (string, error) {
	conn, err := w.Conn()
	if err != nil {
		return stringErr(err.Error())
	}
	defer Close(conn)
	return f(conn)
}

// Strings is a helper function to execute any series of commands
// on a redis.Conn that return a string slice response
func (w *poolClient) Strings(f stringsFunc) ([]string, error) {
	conn, err := w.Conn()
	if err != nil {
		return stringsErr(err.Error())
	}
	defer Close(conn)
	return f(conn)
}

// // DoerCloser Do(es) & Close(s)
// type DoerCloser interface {
// 	Do(string, ...interface{}) (interface{}, error)
// 	Close() error
// }

// // Bool is a helper function to execute any series of commands
// // on a redis.Conn that returns a bool response
// func (w *poolClient) Bool(f DoerCloser, cmd string, args ...interface{}) (bool, error) {
// 	defer f.Close()
// 	return redis.Bool(f.Do(cmd, args...))
// }

// // Int64 is a helper function to execute any series of commands
// // on a redis.Conn that return an int64 response
// func (w *poolClient) Int64(f DoerCloser, cmd string, args ...interface{}) (int64, error) {
// 	defer f.Close()
// 	return redis.Int64(f.Do(cmd, args...))
// }

// // String is a helper function to execute any series of commands
// // on a redis.Conn that return a string response
// func (w *poolClient) String(f DoerCloser, cmd string, args ...interface{}) (string, error) {
// 	defer f.Close()
// 	return redis.String(f.Do(cmd, args...))
// }

// // Strings is a helper function to execute any series of commands
// // on a redis.Conn that return a string slice response
// func (w *poolClient) Strings(f DoerCloser, cmd string, args ...interface{}) ([]string, error) {
// 	defer f.Close()
// 	return redis.Strings(f.Do(cmd, args...))
// }
