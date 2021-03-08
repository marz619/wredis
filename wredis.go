package wredis

import (
	"time"
)

// Wredis Interface
type Wredis interface {
	//
	// Server Commands
	//

	// FlushAll deletes all the keys from all the DB's on the Redis Server.
	//
	// See: http://redis.io/commands/flushall
	FlushAll() error

	// FlushDB deletes all the keys from the configured DB.
	//
	// See: http://redis.io/commands/flushdb
	FlushDB() error

	//
	// Connection Commands
	//

	// Echo echoes the message.
	//
	// See: https://redis.io/commands/echo
	Echo(string) (string, error)

	// Ping returns PONG if no message is provided or the message.
	//
	// See: https://redis.io/commands/ping
	Ping(...string) (string, error)

	// Quit asks the server to close the connection.
	//
	// See: https://redis.io/commands/quit
	Quit() error

	// Select selects the Database specified by the parameter.
	//
	// See: https://redis.io/commands/select
	Select(uint) (Wredis, error)

	// selectable returns if we can call Select on this client
	selectable() bool

	//
	// Keys Commands
	//

	// Del deletes one or more keys from Redis and returns a count of how many keys
	// were actually deleted.
	//
	// See: http://redis.io/commands/del
	Del(...string) (int64, error)

	// Exists checks for the existence of `key` in Redis.
	//
	// See: http://redis.io/commands/exists
	Exists(string) (bool, error)

	// Expire sets a timeout of "seconds" on "key".
	//
	// See: http://redis.io/commands/expire
	Expire(string, int) (bool, error)

	// Keys takes a pattern and returns any/all keys matching the pattern.
	//
	// See: http://redis.io/commands/keys
	Keys(string) ([]string, error)

	// Rename will rename some key "from" to "to".
	//
	// See: `http://redis.io/commands/rename`
	Rename(string, string) error

	//
	// Lists Commands
	//

	// LPush
	// LIndex(string, int64) (string, error)
	LLen(string) (int64, error)
	LPop(string) (string, error)
	LPush(string, ...string) (int64, error)
	RPop(string) (string, error)
	RPush(string, ...string) (int64, error)

	// Sets
	SAdd(string, ...string) (int64, error)
	SCard(string) (int64, error)
	SDiffStore(string, ...string) (int64, error)
	SMembers(string) ([]string, error)
	SUnionStore(string, ...string) (int64, error)

	// Strings
	Append(string, string) (int64, error)
	Get(string) (string, error)
	Incr(string) (int64, error)
	MGet(...string) ([]string, error)
	Set(string, string) error
	SetEx(string, string, uint) error

	// Close
	Close() error

	// Transaction

	// Multi is our entry into Transaction
	//
	// See: https://redis.io/commands/multi
	Multi() (Transaction, error)

	// Watch marks the given keys to be watched fo conditional execution of a
	// transaction.
	//
	// See: https://redis.io/commands/watch
	Watch(...string) error

	// Uwatch flushes all previously watched keys for this transaction.
	//
	// See: https://redis.io/commands/unwatch
	Unwatch() error

	// transacting returns if we're in transaction mode
	transacting() bool

	// convenience functions

	// Exec<type> Funcs
	Bool(boolFunc) (bool, error)
	Int(intFunc) (int, error)
	Int64(int64Func) (int64, error)
	String(stringFunc) (string, error)
	Strings(stringsFunc) ([]string, error)

	// Convenience funcions

	// Appends calls a s
	Appends(string, string, ...string) (int64, error)

	// Delete is an alias for the Del method
	Delete(...string) (int64, error)

	// DelPattern is Del/Delete but uses the provided pattern and executs
	// a Keys command first, to fetch the set of *ALL* keys that match; and
	// subsequently executs a Del(...keys).
	//
	// NOTE: the usage of Keys(pattern) may cause issues in Production
	//		 environments with large databases.
	DelPattern(string) (int64, error)

	// SetExDuration is SetEx but allows a time.Duration to be used as the
	// expiry value. It must be >= 1 * time.Second or it will return an error.
	SetExDuration(string, string, time.Duration) error
}

// Transaction adds the transaction based commands.
//
// NOTE: since Wredis.Multi returns this interface, which embeds a Wredis, you
//       could technically call Multi again in a recursive manner; however our
//       internal implementation will prevent this by returning an error
//
// See: https://redis.io/topics/transactions
type Transaction interface {
	// we can execute all regular Wredis commands within a Transaction.
	Wredis

	// Exec executes all previously queued commands in this transaction.
	//
	// See: https://redis.io/commands/exec
	Exec() ([]interface{}, error)

	// Discard flushes all previously queued commands in this transaction.
	//
	// See: https://redis.io/commands/discard
	Discard() error
}
