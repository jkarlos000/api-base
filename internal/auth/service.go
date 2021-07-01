package auth

import (
	"context"
	"backend/pkg/dbcontext"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"backend/internal/entity"
	"backend/internal/errors"
	"backend/pkg/log"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Service encapsulates the authentication logic.
type Service interface {
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, username, password string) (string, error)
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetName returns the user username.
	GetUsername() string
	// GetRoles returns the user role slice.
	GetRoles() []string
	// HasRole returns the user username.
	HasRole(...string) bool
	// IsUserActive returns the user status
	IsUserActive() bool
}

type service struct {
	db              *dbcontext.DB
	signingKey      string
	tokenExpiration int
	logger          log.Logger
}

// NewService creates a new authentication service.
func NewService(db *dbcontext.DB, signingKey string, tokenExpiration int, logger log.Logger) Service {
	return service{db, signingKey, tokenExpiration, logger}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, username, password string) (string, error) {
	if identity := s.authenticate(ctx, username, password); identity != nil {
		return s.generateJWT(identity)
	}
	return "", errors.Unauthorized("")
}

// authenticate authenticates a user using username and password.
// If username and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, username, password string) Identity {
	logger := s.logger.With(ctx, "user", username)

	user := entity.User{}

	if err := s.db.With(ctx).Select().From("users as u").Where(dbx.HashExp{"u.username": username, "u.is_active": true}).One(&user); err != nil {
		fmt.Println(err)
		return nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println(err)
		logger.Infof("authentication failed")
		return nil
	}

	if err := s.db.With(ctx).Select("r.name as name").
		From("roles as r").
		LeftJoin("role_user as ru", dbx.NewExp("r.id = ru.role_id")).
		Where(dbx.HashExp{"ru.user_id": user.ID}).
		Column(&user.Roles); err != nil {
		fmt.Println(err)
		return nil
	}

	logger.Infof("authentication successful")
	return entity.User{ID: user.GetID(), Username: user.GetUsername(), Roles: user.GetRoles(), IsActive: user.IsActive}
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(identity Identity) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       identity.GetID(),
		"username": identity.GetUsername(),
		"roles":    identity.GetRoles(),
		"status":   identity.IsUserActive(),
		"exp":      time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}
