package nurse

import (
	"context"
	"enfermeria/internal/entity"
	"enfermeria/pkg/dbcontext"
	"enfermeria/pkg/log"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Repository encapsulates the logic to access nurses from the data source.
type Repository interface {
	// Get returns the nurse with the specified nurse ID.
	Get(ctx context.Context, id string) (entity.Nurse, error)
	// Count returns the number of nurses.
	Count(ctx context.Context) (int, error)
	// Query returns the list of nurses with the given offset and limit.
	Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]entity.Nurse, error)
	// Create saves a new nurse in the storage.
	Create(ctx context.Context, user entity.User, nurse entity.Nurse) error
	// Update updates the nurse with given ID in the storage.
	Update(ctx context.Context, user entity.User, nurse entity.Nurse, excludes []string) error
	// Delete removes the nurse with given ID from the storage.
	Delete(ctx context.Context, id string) error
}

// repository persists nurses in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new nurse repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the nurse with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.Nurse, error) {
	var nurse entity.Nurse

	query := fmt.Sprintf("SELECT t.*, u.username, u.first_name, u.last_name, u.is_active FROM telemarketers t LEFT JOIN users u ON u.id = t.user_id WHERE t.id='%v'", id)
	q := r.db.With(ctx).NewQuery(query)

	err := q.One(&nurse)

	return nurse, err
}

// Create saves a new nurse record in the database.
// It returns the ID of the newly inserted nurse record.
func (r repository) Create(ctx context.Context, user entity.User, nurse entity.Nurse) error {
	var role entity.Role

	if err := r.db.With(ctx).Model(&user).Exclude("Roles", "Permissions").Insert(); err != nil {
		return err
	}

	if err := r.db.With(ctx).Model(&nurse).Exclude("FirstName", "LastName", "Username", "IsActive").Insert(); err != nil {
		return err
	}

	if err := r.db.With(ctx).Select().From("roles").Where(dbx.HashExp{"name": "enfermera"}).One(&role); err != nil {
		return err
	}

	if _, err := r.db.With(ctx).Insert("role_user", dbx.Params{
		"role_id": role.ID,
		"user_id": user.ID,
	}).Execute(); err != nil {
		return err
	}

	return nil

}

// Update saves the changes to an nurse in the database.
func (r repository) Update(ctx context.Context, user entity.User, nurse entity.Nurse, excludes []string) error {
	return r.db.With(ctx).Model(&user).Update()
}

// Delete deletes an nurse with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	nurse, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&nurse).Delete()
}

// Count returns the number of the nurse records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("nurses").Row(&count)
	return count, err
}

// Query retrieves the nurse records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]entity.Nurse, error) {
	var nurses []entity.Nurse

	if term == "search" {
		fmt.Println("filters")
		fmt.Println(filters)
		fmt.Println("filters")
		query := fmt.Sprintf("select n.id, u.first_name as first_name, u.last_name as last_name from nurses as n inner join users as u on n.user_id = u.id where CONCAT(u.first_name, ' ', u.last_name) like '%%%s%%'", filters["query"])
		q := r.db.With(ctx).NewQuery(query)
		err := q.All(&nurses)

		return nurses, err
	}

	err := r.db.With(ctx).
		Select("nurses.*", "u.first_name", "u.last_name", "u.username", "u.is_active").
		LeftJoin("users as u", dbx.NewExp("u.id = nurses.user_id")).
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&nurses)
	return nurses, err

}
