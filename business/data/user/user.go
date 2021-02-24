package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// UserRepository ....
type UserRepository struct {
	//log
	db *sqlx.DB
}

// NewRepository ...
func NewRepository(db *sqlx.DB) UserRepository {
	return UserRepository{
		//	log: log,
		db: db,
	}
}

// Query ...
func (u UserRepository) Query(ctx context.Context) ([]User, error) {
	// ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.user.query")
	// defer span.End()

	const q = `
	SELECT
		*
	FROM
		users`
	// ORDER BY
	// 	user_id
	// OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`

	// offset := (pageNumber - 1) * rowsPerPage

	// u.log.Printf("%s: %s: %s", traceID, "user.Query",
	// 	database.Log(q, offset, rowsPerPage),
	// )

	users := []User{}
	if err := u.db.SelectContext(ctx, &users, q); err != nil {
		return nil, err //errors.Wrap(err, "selecting users")
	}

	return users, nil
}
