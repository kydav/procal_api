package routes

import (
	"encoding/json"
	"net/http"
	"procal/entity"
	"procal/repository"
	"procal/services"
	"time"

	"github.com/go-chi/chi"
)

func JournalRoutes(service services.JournalService) func(chi.Router) {
	JournalFunc := func(writer http.ResponseWriter, request *http.Request) {
		JournalHandler(writer, request, service)
	}

	return func(r chi.Router) {
		r.HandleFunc("/journal/{userId}/{date}", JournalFunc)
		r.HandleFunc("/journal", JournalFunc)
	}
}

func JournalHandler(writer http.ResponseWriter, request *http.Request, service services.JournalService) {
	repoInterface := request.Context().Value(repository.ContextKeyRepository)
	sessionRepo, ok := repoInterface.(repository.Repository)
	if !ok {
		panic("could not fetch repo for session")
	}
	JournalRouter(writer, request, services.NewJournalService(sessionRepo.JournalRepository()))
}

func JournalRouter(writer http.ResponseWriter, request *http.Request, service services.JournalService) {
	routePattern := chi.RouteContext(request.Context()).RoutePattern()
	switch request.Method {
	case http.MethodGet:
		if routePattern == "/api/journal/{userId}/{date}" {
			GetUserJournalEntriesByDate(writer, request, service)
		}
	case http.MethodPost:
		if routePattern == "/api/journal" {
			CreateJournalEntry(writer, request, service)
		}
	case http.MethodPut:
		if routePattern == "/api/journal" {
			UpdateJournalEntry(writer, request, service)
		}
	case http.MethodDelete:
		if routePattern == "/api/journal/{id}" {
			DeleteJournalEntry(writer, request, service)
		}
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateJournalEntry(writer http.ResponseWriter, request *http.Request, service services.JournalService) {
	var entry entity.JournalEntry
	if err := json.NewDecoder(request.Body).Decode(&entry); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}
	createdEntry, err := service.CreateEntry(request.Context(), &entry)
	if err != nil {
		http.Error(writer, "Failed to create journal entry", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(createdEntry)
}

func GetUserJournalEntriesByDate(writer http.ResponseWriter, request *http.Request, service services.JournalService) {
	userId := chi.URLParam(request, "userId")
	dateStr := chi.URLParam(request, "date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(writer, "Invalid date format", http.StatusBadRequest)
		return
	}

	entries, err := service.GetEntryByUserAndDate(request.Context(), userId, date)
	if err != nil {
		http.Error(writer, "Journal entry not found", http.StatusNotFound)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(entries)
}

func UpdateJournalEntry(writer http.ResponseWriter, request *http.Request, service services.JournalService) {
	var entry entity.JournalEntry
	if err := json.NewDecoder(request.Body).Decode(&entry); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := service.UpdateEntry(request.Context(), &entry); err != nil {
		http.Error(writer, "Failed to update journal entry", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	json.NewEncoder(writer).Encode(entry)
}

func DeleteJournalEntry(writer http.ResponseWriter, request *http.Request, service services.JournalService) {
	id := chi.URLParam(request, "id")
	if err := service.DeleteEntry(request.Context(), id); err != nil {
		http.Error(writer, "Failed to delete journal entry", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
