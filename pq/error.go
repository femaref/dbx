package pq

import (
    "github.com/lib/pq"
)

func IgnoreErrors(names []string) func(error) error {
    lookup := map[string]bool{}

    for _, name := range names {
        lookup[name] = true
    }

    return func (err error) error {
        pq_err, ok := err.(*pq.Error)

        if ok {
            if _, ok := lookup[pq_err.Code.Name()]; ok {
                err = nil
            }
        }
        return err
    }
}
