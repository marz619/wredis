package wredis

import (
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Config for configuration
type Config struct {
	Cluster         bool
	DB              uint
	Dialer          func(Config) dialFunc
	Host            string
	IdleTimeout     time.Duration
	MaxActive       int
	MaxConnLifetime time.Duration
	MaxIdle         int
	Password        string
	Port            int
	TestOnBorrower  func(Config) borrowFunc
	Wait            bool
	// private config options
	selectable  bool
	transacting bool
}

var nilCfg = Config{}

// Addr return the host:port address from the Config
func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Copy creates copy of the current Config, and modifying it with any provided
// Option(s)
func (c Config) Copy(opts ...Option) (Config, error) {
	// Copy current config
	cfg := Config{
		Cluster:         c.Cluster,
		DB:              c.DB,
		Dialer:          c.Dialer,
		Host:            c.Host,
		IdleTimeout:     c.IdleTimeout,
		MaxActive:       c.MaxActive,
		MaxConnLifetime: c.MaxConnLifetime,
		MaxIdle:         c.MaxIdle,
		Password:        c.Password,
		Port:            c.Port,
		TestOnBorrower:  c.TestOnBorrower,
		Wait:            c.Wait,
		// private config options
		selectable:  c.selectable,
		transacting: c.transacting,
	}

	// apply new/modifying Option(s)
	var err error
	for _, opt := range opts {
		cfg, err = opt(cfg)
		if err != nil {
			return nilCfg, err
		}
	}

	// validate and return
	err = cfg.Validate()
	if err != nil {
		return nilCfg, err
	}

	return cfg, nil
}

// Validate the configuration
func (c Config) Validate() error {
	// cluster -> db/0
	if c.Cluster && c.DB != 0 {
		return errors.New("wredis: cluster supports db/0 only")
	}

	return nil
}

// returns a new default Config
func defaultConfig() Config {
	// sensible? default base values
	cfg := Config{
		Cluster:         false,
		DB:              0,
		Host:            "localhost",
		IdleTimeout:     60 * time.Second,
		MaxActive:       10,
		MaxConnLifetime: time.Hour,
		MaxIdle:         3,
		Port:            6379,
		Wait:            false,
		// private config options
		transacting: false,
		selectable:  true,
	}
	// set the two func config values
	cfg.Dialer = defaultDialer
	cfg.TestOnBorrower = noopTestOnBorrower

	return cfg
}

func newConfig(opts ...Option) (Config, error) {
	// build the configuration using the provided Option(s)
	cfg, err := defaultConfig(), nilErr

	// use Option(s) to modify the default Config
	for _, opt := range opts {
		cfg, err = opt(cfg)
		if err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}

type dialFunc func() (redis.Conn, error)

func defaultDialer(cfg Config) dialFunc {
	return func() (redis.Conn, error) {
		conn, err := redis.Dial("tcp", cfg.Addr(), redis.DialDatabase(int(cfg.DB)))
		if err != nil {
			return nil, err
		}
		// ensure connection is AUTH'd
		if cfg.Password != "" {
			_, err = conn.Do("AUTH", cfg.Password)
			if err != nil {
				return nil, err
			}
		}
		// ensure we're SELECTing the configured DB
		_, err = conn.Do("SELECT", cfg.DB)
		if err != nil {
			return nil, err
		}
		// return this connection
		return conn, nil
	}
}

type borrowFunc func(redis.Conn, time.Time) error

func noopTestOnBorrower(cfg Config) borrowFunc {
	return func(conn redis.Conn, t time.Time) error {
		return nil
	}
}

// Option used to configure the Config
type Option func(Config) (Config, error)

func transacting() Option {
	return func(cfg Config) (Config, error) {
		cfg.transacting = true
		return cfg, nil
	}
}

// unselectable disallows the use of Select. It is a guard against being able
// to call Select against an instance of Wredis returned by a call to Select;
//
// NOTE: this might run into a conflict if the DBSWAP command is used, as the
//       data being returned will be from the swapped database.
func unselectable() Option {
	return func(cfg Config) (Config, error) {
		cfg.selectable = false
		return cfg, nil
	}
}

// Cluster sets the Cluster value in the Config
func Cluster(cluster bool) Option {
	return func(cfg Config) (Config, error) {
		cfg.Cluster = cluster

		// in Cluster mode, disallow Select
		//
		// See: https://redis.io/commands/select
		//
		// When using Redis Cluster, the SELECT command cannot be used, since
		// Redis Cluster only supports database zero.
		if cluster {
			return unselectable()(cfg)
		}

		// otherwise just return the Config
		return cfg, nil
	}
}

// DB sets the DB in the Config
func DB(db uint) Option {
	return func(cfg Config) (Config, error) {
		cfg.DB = db
		return cfg, nil
	}
}

// Dialer sets the Dialer function in the Config
func Dialer(dialer func(Config) dialFunc) Option {
	return func(cfg Config) (Config, error) {
		cfg.Dialer = dialer
		return cfg, nil
	}
}

// Host sets the Host in the Config
func Host(host string) Option {
	return func(cfg Config) (Config, error) {
		if empty(host) {
			return cfg, errors.New("wredis: empty host")
		}
		cfg.Host = host
		return cfg, nil
	}
}

// IdleTimeout sets the IdleTimeout setting in the Config
func IdleTimeout(d time.Duration) Option {
	return func(cfg Config) (Config, error) {
		cfg.IdleTimeout = d
		return cfg, nil
	}
}

// MaxConnLifetime sets the MaxConnLifetime setting in the Config
func MaxConnLifetime(d time.Duration) Option {
	return func(cfg Config) (Config, error) {
		cfg.MaxConnLifetime = d
		return cfg, nil
	}
}

// MaxActive sets the MaxActive setting in the Config
func MaxActive(active int) Option {
	return func(cfg Config) (Config, error) {
		cfg.MaxActive = active
		return cfg, nil
	}
}

// MaxIdle sets the MaxIdle setting in the Config
func MaxIdle(idle int) Option {
	return func(cfg Config) (Config, error) {
		cfg.MaxIdle = idle
		return cfg, nil
	}
}

// Password sets the AUTH password in the Config
func Password(pass string) Option {
	return func(cfg Config) (Config, error) {
		if empty(pass) {
			return cfg, errors.New("wredis: empty password")
		}
		cfg.Password = pass
		return cfg, nil
	}
}

// Port sets the Port in the Config
func Port(port int) Option {
	// disallow reserved ports
	return func(cfg Config) (Config, error) {
		if port <= 1023 {
			return cfg, errors.New("wredis: invalid port")
		}
		cfg.Port = port
		return cfg, nil
	}
}

// TestOnBorrower sets the TestOnBorrower function in the Config
func TestOnBorrower(borrower func(Config) borrowFunc) Option {
	return func(cfg Config) (Config, error) {
		cfg.TestOnBorrower = borrower
		return cfg, nil
	}
}

// Wait sets the Wait in the Config
func Wait(wait bool) Option {
	return func(cfg Config) (Config, error) {
		cfg.Wait = wait
		return cfg, nil
	}
}
