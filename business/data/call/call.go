package call

import (
	"context"
	"database/sql"
	"goapi/foundation/dbase"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// CallRepository ....
type CallRepository struct {
	db *sqlx.DB
}

// NewRepository ...
func NewRepository(db *sqlx.DB) *CallRepository {
	return &CallRepository{
		db: db,
	}
}

// QueryByUser ...
func (crep *CallRepository) QueryByUser(ctx context.Context, uid interface{}) ([]Call, error) {
	const q = `SELECT * FROM calls WHERE "userId" = $1`

	calls := []Call{}
	if err := crep.db.SelectContext(ctx, &calls, q, uid); err != nil {
		return nil, errors.Wrap(err, "CallRepository.QueryByUser:db")
	}

	return calls, nil
}

// CreateForUser ....
func (crep *CallRepository) CreateForUser(ctx context.Context, ccall CreateCall, uid interface{}) (Call, error) {
	const q = `INSERT INTO calls (id, method, path, "userId") VALUES ($1, $2, $3, $4)`

	call := Call{
		ID:     uuid.New().String(),
		Method: ccall.Method,
		Path:   ccall.Path,
	}

	if uid == nil {
		uid = sql.NullString{}
	}

	if _, err := crep.db.ExecContext(ctx, q, call.ID, call.Method, call.Path, uid); err != nil {
		if dbErr, ok := dbase.IsDbError(err); ok {
			return Call{}, dbErr
		}
		return Call{}, errors.Wrap(err, "CallRepository.CreateForUser")
	}
	return call, nil
}
