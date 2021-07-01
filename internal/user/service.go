package user

import (
	"backend/internal/entity"
	"backend/internal/errors"
	"backend/pkg/log"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Service encapsulates usecase logic for users.
type Service interface {
	Get(ctx context.Context, id string) (User, error)
	Me(ctx context.Context) (User, error)
	Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]User, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateUserRequest) (User, error)
	// Update(ctx context.Context, id string, input UpdateUserRequest) (User, error)
	// Delete(ctx context.Context, id string) (User, error)
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

// Me returns the user with the specified context.
func (s service) Me(ctx context.Context) (User, error) {
	user, err := s.repo.Me(ctx)
	if err != nil {
		return User{}, err
	}
	return User{user}, nil
}

// Get returns the user with the specified the user ID.
func (s service) Get(ctx context.Context, id string) (User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return User{}, err
	}
	return User{user}, nil
}

// Create creates a new user.
func (s service) Create(ctx context.Context, req CreateUserRequest) (User, error) {
	if err := req.Validate(); err != nil {
		return User{}, errors.BadRequest(err.Error())
	}

	// identity := auth.CurrentUser(ctx)
	// Validate admin role
	/* if !identity.HasRole("admin") {
		return User{}, errors.BadRequest("Current user has not admin role")
	} */

	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Username = strings.TrimSpace(req.Username)
	if req.Password != nil {
		password := strings.TrimSpace(*req.Password)
		if len(password) > 0 {
			hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.MinCost)
			if err != nil {
				fmt.Println(err)
			}
			password = string(hash)
			req.Password = &password
		}
	}

	id := entity.GenerateID()
	now := time.Now()
	err := s.repo.Create(ctx, entity.User{
		ID:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Password:  *req.Password,
		CreatedAt: now,
		UpdatedAt: &now,
	})
	if err != nil {
		return User{}, err
	}
	return s.Get(ctx, id)
}

// Update updates the user with the specified ID.
func (s service) Update(ctx context.Context, id string, req UpdateUserRequest) (User, error) {
 	if err := req.Validate(); err != nil {
		return User{}, err
	}

	now := time.Now()

	user, err := s.Get(ctx, id)
	if err != nil {
		return user, err
	}
	if user.Username != req.Username {
		user.Username = req.Username
	}
	if user.Email != req.Email {
		user.Email = req.Email
	}

	if user.FirstName != req.FirstName {
		user.FirstName = req.FirstName
	}

	if user.LastName != req.LastName {
		user.LastName = req.LastName
	}

	if req.IsActive != nil && *req.IsActive == false {
		user.IsActive = false
	} else {
		user.IsActive = true
	}

	if req.Password != nil {
		password := strings.TrimSpace(*req.Password)
		if len(password) > 0 {
			hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.MinCost)
			if err != nil {
				fmt.Println(err)
			}
			password = string(hash)
			user.Password = password
			user.UpdatedAt = &now

			if err := s.repo.Update(ctx, user.User); err != nil {
				return user, err
			}
			return user, nil
		}
	}

	user.UpdatedAt = &now

	if err := s.repo.Update(ctx, user.User); err != nil {
		return user, err
	}
	return user, nil
}

// Delete deletes the user with the specified ID.
// func (s service) Delete(ctx context.Context, id string) (User, error) {
// 	user, err := s.Get(ctx, id)
// 	if err != nil {
// 		return User{}, err
// 	}
// 	if err = s.repo.Delete(ctx, id); err != nil {
// 		return User{}, err
// 	}
// 	return user, nil
// }

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
