package user

import (
	"backend/pkg/log"
	"backend/pkg/pagination"
	"encoding/json"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}
	r.Use(authHandler)
	// the following endpoints require a valid JWT
	r.Get("/users/<id>", res.get)
	r.Get("/users", res.query)
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
