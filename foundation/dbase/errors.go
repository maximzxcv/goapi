package dbase

import (
	"errors"

	"github.com/lib/pq"
)

var (
	ErrNotExist     = errors.New("not exist")
	ErrAlreadyExist = errors.New("already exist")
)

// IsDbError check if error came from DB and wraps it to supported one
func IsDbError(err error) (error, bool) {
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code.Name() == "unique_violation" {
			return ErrAlreadyExist, true
		}
	}
	return err, false
}
