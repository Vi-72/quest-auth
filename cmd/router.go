package cmd

import (
	"errors"
	stdhttp "net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	openapihttp "github.com/Vi-72/quest-auth/api/http/auth/v1"
	httpmiddleware "github.com/Vi-72/quest-auth/internal/adapters/in/http/middleware"
	"github.com/Vi-72/quest-auth/internal/adapters/in/http/problems"
	"github.com/Vi-72/quest-auth/internal/pkg/errs"
)

const apiV1Prefix = "/api/v1"

func NewRouter(root *CompositionRoot) stdhttp.Handler {
	router := chi.NewRouter()

	// --- Базовые middleware ---
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Logger)

	// Load OpenAPI spec
	swagger, err := openapihttp.GetSwagger()
	if err != nil {
		panic("failed to load OpenAPI spec: " + err.Error())
	}

	swagger.Servers = []*openapi3.Server{{URL: apiV1Prefix}}

	// --- Health check ---
	router.Get("/health", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Swagger JSON
	router.Get("/openapi.json", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		spec, err := openapihttp.GetSwagger()
		if err != nil {
			stdhttp.Error(w, "failed to load OpenAPI spec: "+err.Error(), stdhttp.StatusInternalServerError)
			return
		}
		bytes, err := spec.MarshalJSON()
		if err != nil {
			stdhttp.Error(w, "failed to marshal OpenAPI: "+err.Error(), stdhttp.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write(bytes)
	})

	// Swagger UI
	router.Get("/docs", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
			  <meta charset="UTF-8">
			  <title>Swagger UI</title>
			  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
			</head>
			<body>
			  <div id="swagger-ui"></div>
			  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
			  <script>
				window.onload = () => {
				  SwaggerUIBundle({
					url: "/openapi.json",
					dom_id: "#swagger-ui",
				  });
				};
			  </script>
			</body>
			</html>
		`))
	})

	strictHandler := root.NewAPIHandler()

	// Create StrictHandler with custom error handling for parameter parsing and validation
	apiHandler := openapihttp.NewStrictHandlerWithOptions(strictHandler, nil, openapihttp.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
			// Handle parameter parsing errors with detailed messages
			problem := problems.NewBadRequest("Invalid request parameters: " + err.Error())
			problem.WriteResponse(w)
		},
		ResponseErrorHandlerFunc: func(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
			// Check if it's a domain validation error from application layer
			var domainValidationErr *errs.DomainValidationError
			if errors.As(err, &domainValidationErr) {
				// Convert to 400 Bad Request
				problem := problems.NewBadRequest(domainValidationErr.Error())
				problem.WriteResponse(w)
				return
			}

			// Check if it's a not found error from application layer
			var notFoundErr *errs.NotFoundError
			if errors.As(err, &notFoundErr) {
				// Convert to 404 Not Found
				problem := problems.NewNotFound(notFoundErr.Error())
				problem.WriteResponse(w)
				return
			}

			// Handle other response errors
			problem := problems.NewBadRequest("Response error: " + err.Error())
			problem.WriteResponse(w)
		},
	})

	// Create API router with OpenAPI validation middleware
	apiRouter := chi.NewRouter()

	// Add OpenAPI validation middleware
	validationMW, err := httpmiddleware.NewOpenAPIValidationMiddleware(swagger)
	if err == nil {
		apiRouter.Use(validationMW.Validate)
	}

	openapihttp.HandlerFromMux(apiHandler, apiRouter)

	router.Mount(apiV1Prefix, apiRouter)

	return router
}
