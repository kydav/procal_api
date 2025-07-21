package repository

import (
	"context"
	"net/http"
)

func BuildRepositoryWithContextMiddlware(repo Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			// For a given request, the repository context can outlive the request in the case of background / async
			// processing. Allow it to process without the cancellation context of the request.
			contextWithoutCancel := context.WithoutCancel(request.Context())
			contextRepo := repo.CreateRepositoryWithContext(contextWithoutCancel)
			updatedContext := context.WithValue(request.Context(), "Repository", contextRepo)
			request = request.WithContext(updatedContext)
			next.ServeHTTP(writer, request)
		})
	}
}
