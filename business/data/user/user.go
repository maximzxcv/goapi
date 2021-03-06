package user

import (
	"context"
	"database/sql"
	"goapi/business/auth"
	"goapi/foundation/dbase"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository ....
type UserRepository struct {
	db *sqlx.DB
}

// NewRepository ...
func NewRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Query ...
func (urep UserRepository) Query(ctx context.Context) ([]User, error) {
	const q = `SELECT * FROM users`

	usrs := []User{}
	if err := urep.db.SelectContext(ctx, &usrs, q); err != nil {
		return nil, errors.Wrap(err, "Query:db")
	}

	return usrs, nil
}

// QueryByID .....
func (urep UserRepository) QueryByID(ctx context.Context, uid string) (User, error) {
	const q = `SELECT * FROM users AS u WHERE u.id=$1`
	var usr User
	if err := urep.db.GetContext(ctx, &usr, q, uid); err != nil {
		if err == sql.ErrNoRows {
			return usr, dbase.ErrNotExist
		}
		return usr, errors.Wrap(err, "QueryByID:db")
	}
	return usr, nil
}

// Delete ...
func (urep UserRepository) Delete(ctx context.Context, uid string) error {
	const q = `DELETE FROM users AS u WHERE u.id=$1`
	if _, err := urep.db.ExecContext(ctx, q, uid); err != nil {
		return errors.Wrap(err, "Delete:db")
	}
	return nil
}

// Create ....
func (urep UserRepository) Create(ctx context.Context, cusr CreateUser) (User, error) {
	const q = `INSERT INTO users (id, name, password) VALUES ($1, $2, $3)`

	hash, err := bcrypt.GenerateFromPassword([]byte(cusr.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.Wrap(err, "Create:Encrypt password")
	}

	usr := User{
		ID:   uuid.New().String(),
		Name: cusr.Name,
	}

	if _, err := urep.db.ExecContext(ctx, q, usr.ID, usr.Name, hash); err != nil {
		if dbErr, ok := dbase.IsDbError(err); ok {
			return User{}, dbErr
		}
		return User{}, errors.Wrap(err, "UserRepository.Create")
	}
	return usr, nil
}

// Update ....
func (urep UserRepository) Update(ctx context.Context, uid string, uusr UpdateUser) (User, error) {
	const q = `UPDATE users SET
	 	name=$2, 
	 	password = $3
	 WHERE id=$1`

	usr, err := urep.QueryByID(ctx, uid)
	if err != nil {
		return User{}, errors.Wrap(err, "Update:cannot get user by id")
	}

	if uusr.Name != "" {
		usr.Name = uusr.Name
	}

	var hash []byte
	if uusr.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(uusr.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, errors.Wrap(err, "Update:Encrypt password")
		}
		usr.Password = hash
	}

	if _, err := urep.db.ExecContext(ctx, q, usr.ID, usr.Name, hash); err != nil {
		if dbErr, ok := dbase.IsDbError(err); ok {
			return User{}, dbErr
		}
		return User{}, errors.Wrap(err, "Update:db")
	}

	return usr, nil
}

// CheckAuth validates if password is correct for user
func (urep UserRepository) CheckAuth(ctx context.Context, username string, password string) (User, error) {
	const q = `SELECT * FROM users AS u WHERE u.name=$1`
	var usr User
	if err := urep.db.GetContext(ctx, &usr, q, username); err != nil {
		if err == sql.ErrNoRows {
			return usr, auth.ErrNotAuthorised
		}
		return usr, errors.Wrap(err, "CheckAuth:db")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)); err != nil {
		return User{}, auth.ErrNotAuthorised
	}

	return usr, nil
}
