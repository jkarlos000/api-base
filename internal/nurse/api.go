package nurse

import (
	"encoding/json"
	"enfermeria/internal/errors"
	"enfermeria/pkg/log"
	"enfermeria/pkg/pagination"
	"github.com/go-ozzo/ozzo-routing/v2"
	"net/http"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/nurses/<id>", res.get)
	r.Get("/nurses", res.query)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Post("/nurses", res.create)
	r.Put("/nurses/<id>", res.update)
	r.Delete("/nurses/<id>", res.delete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	nurse, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(nurse)
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
	nurses, err := r.service.Query(ctx, pages.Offset(), pages.Limit(), term, filters)
	if err != nil {
		return err
	}
	pages.Items = nurses
	return c.Write(pages)
}

func (r resource) create(c *routing.Context) error {
	var input CreateNurseRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	nurse, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(nurse, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateNurseRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	nurse, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(nurse)
}

func (r resource) delete(c *routing.Context) error {
	nurse, err := r.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(nurse)
}
