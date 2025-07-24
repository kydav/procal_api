package repository

import (
	"context"
	"net/http"
)

func BuildRepositoryWithContextMiddlware(repo Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			contextWithoutCancel := context.WithoutCancel(request.Context())
			contextRepo := repo.CreateRepositoryWithContext(contextWithoutCancel)
			updatedContext := context.WithValue(request.Context(), ContextKeyRepository, contextRepo)
			request = request.WithContext(updatedContext)
			next.ServeHTTP(writer, request)
		})
	}
}

type ContextKey string

const (
	ContextKeyRepository ContextKey = "Repository"
)
