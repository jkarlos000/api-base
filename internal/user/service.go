package user

import (
	"backend/internal/entity"
	"backend/pkg/log"
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
)

// Service encapsulates usecase logic for users.
type Service interface {
	Get(ctx context.Context, id string) (User, error)
	Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]User, error)
	Count(ctx context.Context) (int, error)

}

// User represents the data about an user.
type User struct {
	entity.User
}

// CreateUserRequest represents an user creation request.
type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  *string `json:"password"`
	Email		string `json:"email"`
}

// Validate validates the CreateUserRequest fields.
func (m CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FirstName, validation.Required, validation.Length(3, 50), validation.Match(regexp.MustCompile("^([A-Za-z']+ )+[A-Za-z']+$|^[A-Za-z']+$"))),
		validation.Field(&m.LastName, validation.Required, validation.Length(3, 50), validation.Match(regexp.MustCompile("^([A-Za-z']+ )+[A-Za-z']+$|^[A-Za-z']+$"))),
		validation.Field(&m.Username, validation.Required, validation.Length(3, 50), validation.Match(regexp.MustCompile("^([0-9A-Za-z]+ )+[0-9A-Za-z]+$|^[0-9A-Za-z]+$"))),
		validation.Field(&m.Password, validation.Required, validation.Length(0, 150)),
		validation.Field(&m.Email, validation.Required, validation.Length(0, 150)),
	)
}

// UpdateUserRequest represents an user update request.
type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  *string `json:"password"`
	Email		string `json:"email"`
	IsActive	*bool `json:"is_active"`
}

// Validate validates the CreateUserRequest fields.
func (m UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FirstName, validation.Required, validation.Length(3, 50), validation.Match(regexp.MustCompile("^([A-Za-z']+ )+[A-Za-z']+$|^[A-Za-z']+$"))),
		validation.Field(&m.LastName, validation.Required, validation.Length(3, 50), validation.Match(regexp.MustCompile("^([A-Za-z']+ )+[A-Za-z']+$|^[A-Za-z']+$"))),
		validation.Field(&m.Username, validation.Required, validation.Length(3, 50), validation.Match(regexp.MustCompile("^([0-9A-Za-z]+ )+[0-9A-Za-z]+$|^[0-9A-Za-z]+$"))),
		// validation.Field(&m.Password, validation.Length(0, 150)),
		validation.Field(&m.Email, validation.Required, validation.Length(0, 150)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new user service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the user with the specified the user ID.
func (s service) Get(ctx context.Context, id string) (User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return User{}, err
	}
	return User{user}, nil
}

// Count returns the number of users.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// Query returns the users with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]User, error) {
	items, err := s.repo.Query(ctx, offset, limit, term, filters)
	if err != nil {
		return nil, err
	}
	result := []User{}
	for _, item := range items {
		result = append(result, User{item})
	}
	return result, nil
}
