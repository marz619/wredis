package wredis

// See: https://redis.io/topics/transactions

func (w *impl) Multi() (Transaction, error) {
	return nil, nil
}

func (w *impl) Discard() error {
	return nil
}

func (w *impl) Exec() ([]interface{}, error) {
	return nil, nil
}

func (w *impl) Watch(keys ...string) error {
	return nil
}

func (w *impl) Unwatch() error {
	return nil
}

func (w *impl) transacting() bool {
	return w.cfg.transacting
}
