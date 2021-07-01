package session

import (
	"context"
	"backend/internal/entity"
	"backend/pkg/log"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// Service encapsulates uses case logic for sessions.
type Service interface {
	Get(ctx context.Context, id string) (Session, error)
	Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]Session, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateSessionRequest) (Session, error)
	Update(ctx context.Context, id string, input UpdateSessionRequest) (Session, error)
	Delete(ctx context.Context, id string) (Session, error)
}

// Session represents the data about an session.
type Session struct {
	entity.Session
}

// CreateSessionRequest represents an user creation request.
type CreateSessionRequest struct {
	FirstName            string  `json:"first_name"`
	LastName             string  `json:"last_name"`
	Username             string  `json:"username"`
	Password             string  `json:"password"`
	Latitude 			*float32 `json:"latitude"`
	Longitude 			*float32 `json:"longitude"`
}

// Validate validates the CreateSessionRequest fields.
func (m CreateSessionRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FirstName, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.LastName, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Username, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Password, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Latitude, validation.Required),
		validation.Field(&m.Longitude, validation.Required),
	)
}
// UpdateSessionRequest represents an user update request.
type UpdateSessionRequest struct {
	FirstName            string  `json:"first_name"`
	LastName             string  `json:"last_name"`
	Username             string  `json:"username"`
	Password             *string `json:"password"`
	Latitude 			*float32 `json:"latitude"`
	Longitude 			*float32 `json:"longitude"`
	Rol					entity.Role `json:"rol"`
	IsActive			*bool `json:"is_active"`
}

// Validate validates the CreateUserRequest fields.
func (m UpdateSessionRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FirstName, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.LastName, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Username, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Latitude, validation.Required),
		validation.Field(&m.Longitude, validation.Required),
		validation.Field(&m.Rol, validation.Required),
		validation.Field(&m.IsActive, validation.Required),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new session service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the session with the specified the session ID.
func (s service) Get(ctx context.Context, id string) (Session, error) {
	session, err := s.repo.Get(ctx, id)
	if err != nil {
		return Session{}, err
	}
	return Session{session}, nil
}

// Create creates a new session.
func (s service) Create(ctx context.Context, req CreateSessionRequest) (Session, error) {
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Username = strings.TrimSpace(req.Username)

	sessionID := entity.GenerateID()
	userID := entity.GenerateID()

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return Session{}, err
	}

	if err := req.Validate(); err != nil {
		return Session{}, err
	}
	now := time.Now()
	err = s.repo.Create(ctx, entity.User{
		ID:        userID,
		Username:  req.Username,
		Password:  string(hash),
		IsActive:  true,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: now,
		UpdatedAt: &now,
	}, entity.Session{
		ID:             sessionID,
		UserID:         userID,
		HaveFullAccess: false,
		IsActive:       true,
		IsWorking:      false,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
	})
	if err != nil {
		return Session{}, err
	}
	return s.Get(ctx, sessionID)
}

// Update updates the session with the specified ID.
func (s service) Update(ctx context.Context, id string, req UpdateSessionRequest) (Session, error) {
	var user entity.User
	excludes := []string{"Roles", "Permissions", "CreatedAt"}

	if err := req.Validate(); err != nil {
		return Session{}, err
	}

	now := time.Now()

	res, err := s.Get(ctx, id)
	if err != nil {
		return res, err
	}
	session := res.Session
	session.Latitude = req.Latitude
	session.Longitude = req.Longitude

	user.ID = session.UserID
	user.Username = req.Username
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.UpdatedAt = &now

	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.MinCost)
		if err != nil {
			fmt.Println(err)
		}

		user.Password = string(hash)
	} else {
		excludes = append(excludes, "Password")
	}

	if req.IsActive != nil && *req.IsActive == false {
		user.IsActive = false
		res.IsActive = false
	} else {
		user.IsActive = true
		res.IsActive = true
	}
	// Fixme refactor repo.Update - remove user and excludes
	if err := s.repo.Update(ctx, user, session, excludes); err != nil {
		return res, err
	}

	res.Username = req.Username
	res.FirstName = req.FirstName
	res.LastName = req.LastName
	res.Latitude = req.Latitude
	res.Longitude = req.Longitude

	return res, nil
}

// Delete deletes the session with the specified ID.
func (s service) Delete(ctx context.Context, id string) (Session, error) {
	session, err := s.Get(ctx, id)
	if err != nil {
		return Session{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Session{}, err
	}
	return session, nil
}

// Count returns the number of sessions.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// Query returns the sessions with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]Session, error) {
	items, err := s.repo.Query(ctx, offset, limit, term, filters)
	if err != nil {
		return nil, err
	}
	result := []Session{}
	for _, item := range items {
		result = append(result, Session{item})
	}
	return result, nil
}
