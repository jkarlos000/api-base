package session

import (
	"backend/internal/auth"
	"backend/internal/entity"
	"backend/pkg/log"
	"context"
	"encoding/hex"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
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
	Tittle		string `json:"tittle"`
	Description	string `json:"description"`
	Password	*string `json:"password"`
}

// Validate validates the CreateSessionRequest fields.
func (m CreateSessionRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Tittle, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Description, validation.Length(0, 128)),
		validation.Field(&m.Password, validation.Length(0, 128)),
	)
}
// UpdateSessionRequest represents an user update request.
type UpdateSessionRequest struct {
	ID			string `json:"id"`
	Tittle		string `json:"tittle"`
	Description	string `json:"description"`
	Password	*string `json:"password"`
	IsActive	*bool `json:"is_active"`
}

// Validate validates the CreateUserRequest fields.
func (m UpdateSessionRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ID, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Tittle, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Description, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Password, validation.Length(0, 128)),
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
	req.Tittle = strings.TrimSpace(req.Tittle)
	req.Description = strings.TrimSpace(req.Description)
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

	sessionID := entity.GenerateID()
	userID := auth.CurrentUser(ctx).GetID()
	slugName := generateSecureSlug(5)
	if err := req.Validate(); err != nil {
		return Session{}, err
	}
	now := time.Now()
	err := s.repo.Create(ctx, entity.Session{
		ID:          sessionID,
		Owner:       userID,
		Tittle:      req.Tittle,
		Description: req.Description,
		Slug:        slugName,
		Password:    *req.Password,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   &now,
	})
	if err != nil {
		return Session{}, err
	}
	return s.Get(ctx, sessionID)
}

// Update updates the session with the specified ID.
func (s service) Update(ctx context.Context, id string, req UpdateSessionRequest) (Session, error) {
	req.Tittle = strings.TrimSpace(req.Tittle)
	req.Description = strings.TrimSpace(req.Description)
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
	if err := req.Validate(); err != nil {
		return Session{}, err
	}

	res, err := s.Get(ctx, id)
	if err != nil {
		return res, err
	}
	session := res.Session
	if session.Tittle != req.Tittle {
		session.Tittle = req.Tittle
	}
	if session.Description != req.Description {
		session.Description = req.Description
	}
	if req.IsActive != nil && *req.IsActive == false {
		session.IsActive = false
	} else {
		session.IsActive = true
	}
	if req.Password != nil {
		if len(*req.Password) > 0 {
			session.Password = *req.Password
		}
	}
	now := time.Now()
	session.UpdatedAt = &now

	// Fixme refactor repo.Update - remove user and excludes
	if err := s.repo.Update(ctx, session, []string{}); err != nil {
		return res, err
	}

	q, _ := s.Get(ctx, req.ID)

	return q, nil
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

func generateSecureSlug(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
