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

func intErr(msg string) (int, error) {
	return 0, errors.New(msg)
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
	return fmt.Errorf("wredis: %s requires unsafe impl. See wredis.Unsafe", method)
}

//
// utilties
//

// returns true if a string whos spaces are trimmed is empty
func empty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// returns true if any string in the slice sastisfies the predicate
func any(ss []string, pred func(string) bool) bool {
	for _, s := range ss {
		if pred(s) {
			return true
		}
	}
	return false
}
