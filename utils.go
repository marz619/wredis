package wredis

import (
	"errors"
	"fmt"
	"strings"
)

//
// error helper functions
//

func boolErr(msg string) (bool, error) {
	return false, errors.New(msg)
}

func int64Err(msg string) (int64, error) {
	return 0, errors.New(msg)
}

func stringErr(msg string) (string, error) {
	return "", errors.New(msg)
}

func stringsErr(msg string) ([]string, error) {
	return nil, errors.New(msg)
}

func unsafeErr(method string) error {
	return fmt.Errorf("wredis: %s requires unsafe poolClient. See wredis.Unsafe", method)
}

//
// utilties
//

func empty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func any(ss []string, pred func(string) bool) bool {
	for _, s := range ss {
		if pred(s) {
			return true
		}
	}
	return false
}
