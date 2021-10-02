package user

import (
	"backend/internal/errors"
	"backend/pkg/log"
	"backend/pkg/pagination"
	"encoding/json"
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Post("/users", res.create)
	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Get("/users/me", res.me)
	r.Get("/users/<id>", res.get)
	r.Get("/users", res.query)
	// r.Put("/users/<id>", res.update)
	// r.Delete("/users/<id>", res.delete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	user, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(user)
}

func (r resource) me(c *routing.Context) error {
	user, err := r.service.Me(c.Request.Context())
	if err != nil {
		return err
	}

	return c.Write(user)
}

func (r resource) query(c *routing.Context) error {
	term := c.Query("term")
	filters := make(map[string]interface{})

	// convert JSON string filters to map
	_ = json.Unmarshal([]byte(c.Query("filters")), &filters)

	ctx := c.Request.Context()
	count, err := r.service.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)
	users, err := r.service.Query(ctx, pages.Offset(), pages.Limit(), term, filters)
	if err != nil {
		return err
	}
	pages.Items = users
	return c.Write(pages)
}


func (r resource) create(c *routing.Context) error {
	var input CreateUserRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	user, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(user, http.StatusCreated)
}

// func (r resource) update(c *routing.Context) error {
// 	var input UpdateUserRequest
// 	if err := c.Read(&input); err != nil {
// 		r.logger.With(c.Request.Context()).Info(err)
// 		return errors.BadRequest("")
// 	}

// 	user, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Write(user)
// }

// func (r resource) delete(c *routing.Context) error {
// 	user, err := r.service.Delete(c.Request.Context(), c.Param("id"))
// 	if err != nil {
// 		return err
// 	}

// 	return c.Write(user)
// }
