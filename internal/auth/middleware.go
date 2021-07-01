package auth

import (
	"context"
	"backend/internal/entity"
	"backend/internal/errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/auth"
	"net/http"
)

// Handler returns a JWT-based authentication middleware.
func Handler(verificationKey string) routing.Handler {
	return auth.JWT(verificationKey, auth.JWTOptions{TokenHandler: handleToken})
}

// handleToken stores the user identity in the request context so that it can be accessed elsewhere.
func handleToken(c *routing.Context, token *jwt.Token) error {
	var roles []string
	interfaceRoles := token.Claims.(jwt.MapClaims)["roles"].([]interface{})

	for _, v := range interfaceRoles {
		roles = append(roles, fmt.Sprint(v))
	}

	ctx := WithUser(
		c.Request.Context(),
		token.Claims.(jwt.MapClaims)["id"].(string),
		token.Claims.(jwt.MapClaims)["username"].(string),
		roles,
		token.Claims.(jwt.MapClaims)["status"].(bool),
	)

	c.Request = c.Request.WithContext(ctx)
	return nil
}

type contextKey int

const (
	userKey contextKey = iota
)

// WithUser returns a context that contains the user identity from the given JWT.
func WithUser(ctx context.Context, id, username string, roles []string, isActive bool) context.Context {
	return context.WithValue(ctx, userKey, entity.User{ID: id, Username: username, Roles: roles, IsActive: isActive})
}

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
func CurrentUser(ctx context.Context) Identity {
	if user, ok := ctx.Value(userKey).(entity.User); ok {
		return user
	}

	return nil
}

// MockAuthHandler creates a mock authentication middleware for testing purpose.
// If the request contains an Authorization header whose value is "TEST", then
// it considers the user is authenticated as "Tester" whose ID is "100".
// It fails the authentication otherwise.
func MockAuthHandler(c *routing.Context) error {
	if c.Request.Header.Get("Authorization") != "TEST" {
		return errors.Unauthorized("")
	}
	ctx := WithUser(c.Request.Context(), "100", "Tester", []string{"admin"}, true)
	c.Request = c.Request.WithContext(ctx)
	return nil
}

// MockAuthHeader returns an HTTP header that can pass the authentication check by MockAuthHandler.
func MockAuthHeader() http.Header {
	header := http.Header{}
	header.Add("Authorization", "TEST")
	return header
}
