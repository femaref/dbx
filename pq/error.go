package pq

import (
	"github.com/lib/pq"
)

func IgnoreErrors(names ...string) func(error) error {
	fn := CheckErrors(names...)
	return func(err error) error {
		if fn(err) {
			return nil
		}
		return err
	}
}

func CheckErrors(names ...string) func(error) bool {
	lookup := map[string]bool{}

	for _, name := range names {
		lookup[name] = true
	}

	return func(err error) bool {
		if err == nil {
			return false
		}
		pq_err, ok := err.(*pq.Error)

		if ok {
			if ok := lookup[pq_err.Code.Name()]; ok {
				return true
			}
		}
		return false
	}
}

var SkipUniqueConstraint = IgnoreErrors("unique_violation")
var CheckUniqueConstraint = CheckErrors("unique_violation")
