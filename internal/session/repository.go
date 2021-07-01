package session

import (
	"backend/internal/entity"
	"backend/pkg/dbcontext"
	"backend/pkg/log"
	"context"
	"encoding/hex"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"math/rand"
)

// Repository encapsulates the logic to access sessions from the data source.
type Repository interface {
	// Get returns the session with the specified session ID.
	Get(ctx context.Context, id string) (entity.Session, error)
	// Count returns the number of sessions.
	Count(ctx context.Context) (int, error)
	// Query returns the list of sessions with the given offset and limit.
	Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]entity.Session, error)
	// Create saves a new session in the storage.
	Create(ctx context.Context, session entity.Session) error
	// Update updates the session with given ID in the storage.
	Update(ctx context.Context, session entity.Session, excludes []string) error
	// Delete removes the session with given ID from the storage.
	Delete(ctx context.Context, id string) error
}

// repository persists sessions in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new session repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the session with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.Session, error) {
	var session entity.Session

	query := fmt.Sprintf("SELECT t.*, u.username, u.first_name, u.last_name, u.is_active FROM telemarketers t LEFT JOIN users u ON u.id = t.user_id WHERE t.id='%v'", id)
	q := r.db.With(ctx).NewQuery(query)

	err := q.One(&session)

	return session, err
}

// Create saves a new session record in the database.
// It returns the ID of the newly inserted session record.
func (r repository) Create(ctx context.Context, session entity.Session) error {
	var role entity.Role
	if err := r.db.With(ctx).Select().From("roles").Where(dbx.HashExp{"name": "owner"}).One(&role); err != nil {
		return err
	}
	var q entity.Session
	b := true
	for b {
		if err := r.db.With(ctx).Select().From("sessions").Where(dbx.HashExp{"slug": session.Slug}).One(&q); err == nil {
			session.Slug = func(length int) string {
				b := make([]byte, length)
				if _, err := rand.Read(b); err != nil {
					return ""
				}
				return hex.EncodeToString(b)
			}(5)
		} else {
			b = false
		}
	}

	if _, err := r.db.With(ctx).Insert("rooms", dbx.Params{
		"session_id": session.ID,
		"user_id": session.Owner,
		"role_id": role.ID,
	}).Execute(); err != nil {
		return err
	}

	if err := r.db.With(ctx).Model(&session).Insert(); err != nil {
		return err
	}

	return nil

}

// Update saves the changes to an session in the database.
func (r repository) Update(ctx context.Context, session entity.Session, excludes []string) error {
	return r.db.With(ctx).Model(&session).Update()
}

// Delete deletes an session with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	session, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&session).Delete()
}

// Count returns the number of the session records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("sessions").Row(&count)
	return count, err
}

// Query retrieves the session records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]entity.Session, error) {
	var sessions []entity.Session

	if term == "user" {
		query := fmt.Sprintf("select s.id, u.id as owner, u.first_name as first_name, u.last_name as last_name, s.slug, s.updated_at from sessions as s inner join users as u on s.owner = u.id where s.owner = %%%s%%", filters["query"])
		q := r.db.With(ctx).NewQuery(query)
		err := q.All(&sessions)

		return sessions, err
	}

	if term == "active" {
		active := 1
		if filters["query"] == "active" {
			active = 1
		} else if filters["query"] == "inactive" {
			active = 0
		}
		query := fmt.Sprintf("select s.id, u.id as owner, u.first_name as first_name, u.last_name as last_name, s.slug, s.updated_at from sessions as s inner join users as u on s.owner = u.id where s.is_active = %%%d%%", active)
		q := r.db.With(ctx).NewQuery(query)
		err := q.All(&sessions)

		return sessions, err
	}

	if term == "active" {
		query := fmt.Sprintf("select s.id, u.id as owner, u.first_name as first_name, u.last_name as last_name, s.slug, s.updated_at from sessions as s inner join users as u on s.owner = u.id where CONCAT(u.first_name, ' ', u.last_name) like '%%%s%%'", filters["query"])
		q := r.db.With(ctx).NewQuery(query)
		err := q.All(&sessions)

		return sessions, err
	}

	err := r.db.With(ctx).
		Select("sessions.*", "u.first_name", "u.last_name", "u.username", "u.is_active").
		LeftJoin("users as u", dbx.NewExp("u.id = sessions.owner")).
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&sessions)
	return sessions, err

}
