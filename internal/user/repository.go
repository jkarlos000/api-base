package user

import (
	"backend/internal/entity"
	"backend/internal/errors"
	"backend/pkg/dbcontext"
	"backend/pkg/log"
	"context"
	"fmt"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	// Get returns the user with the specified user ID.
	Get(ctx context.Context, id string) (entity.User, error)
	// Count returns the number of users.
	Count(ctx context.Context) (int, error)
	// Query returns the list of users with the given offset and limit.
	Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]entity.User, error)
	// Create saves a new user in the storage.
	Create(ctx context.Context, user entity.User) error
}

// repository persists users in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new user repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the user with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Select().Model(id, &user)
	return user, err
}

// Update saves the changes to an user in the database.
func (r repository) Update(ctx context.Context, user entity.User) error {
	return r.db.With(ctx).Model(&user).Update()
}

// Count returns the number of the user records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("users").Row(&count)
	return count, err
}

// Query retrieves the user records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]entity.User, error) {
	var users []entity.User
	if term == "search" {
		switch filters["role"] {
		case "owner":
			query := fmt.Sprintf("%%%v%%", filters["query"])
			err := r.db.With(ctx).
				Select("users.first_name", "users.last_name", "users.id as id").
				InnerJoin("rooms as r", dbx.NewExp("r.user_id = users.id")).
				Where(dbx.NewExp("CONCAT(users.first_name, ' ', users.last_name) like {:query}", dbx.Params{"query": query})).
				OrderBy("users.id").
				Offset(int64(offset)).
				Limit(int64(limit)).
				All(&users)
			return users, err
		}
	}

	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&users)
	return users, err
}
// Create saves a new user record in the database.
// It returns the ID of the newly inserted user record.
func (r repository) Create(ctx context.Context, user entity.User) error {
	var count int

	if err := r.db.With(ctx).Select("COUNT(*)").From("users").Where(dbx.HashExp{"username": user.Username}).Row(&count); err != nil {
		return err
	}

	if count > 0 {
		return errors.BadRequest("username already exists")
	}

	user.IsActive = true

	return r.db.With(ctx).Model(&user).Insert()
}
