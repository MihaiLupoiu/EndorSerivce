package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/core/services"
)

type HandlerHTTP struct {
	svc *services.EndorService
	v   *validator.Validate
	r   *chi.Mux
}

// configureRoutes configures the routes for the HTTP handler.
func (h *HandlerHTTP) configureRoutes() {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	// TODO: Change to specific allowedOrigins e.g: localhost:8080 to use Access-Control-Allow-Credentials
	h.r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	h.r.Use(middleware.RequestID)
	h.r.Use(middleware.Logger)

	// Attack HTTP handlers
	h.r.Route("/attack", func(r chi.Router) {
		r.Group(func(r chi.Router) { // Group use to apply middleweres only to this path
			r.Use(render.SetContentType(render.ContentTypeJSON))
			r.Post("/", h.getTarget)
		})
	})
}

// getTarget is the HTTP handler for the "/attack" endpoint.
func (h *HandlerHTTP) getTarget(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := &AttackRequest{}
	if err := render.DecodeJSON(r.Body, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest))
		return
	}

	// Validate data
	// We can implement the validation as a middleware and add to the router.
	// This would require more work as we need to make a copy of the request body
	// and the validation middlewere shuold be reuseble by different routes.
	err := h.validateHTTPAttackPOST(data)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest))
		return
	}

	// Convert data to Domain Model
	attackData, err := data.ConvertToAttackDataModel()
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest))
		return
	}

	// Return
	target, err := h.svc.Attack(attackData)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest))
		return
	}

	var res = &AttackReportResponse{
		Casualties: target.Casualties,
		Generation: target.Generation,
		Target: &Coordinate{
			X: &target.Target.X,
			Y: &target.Target.Y,
		},
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// validateHTTPAttackPOST validates the HTTP POST request data for the "/attack" endpoint.
func (h *HandlerHTTP) validateHTTPAttackPOST(data *AttackRequest) error {
	err := h.v.Struct(data)
	if err != nil {
		return err
	}

	return nil
}
