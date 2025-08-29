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

func MealRoutes() func(chi.Router) {
	MealFunc := func(writer http.ResponseWriter, request *http.Request) {
		MealHandler(writer, request)
	}

	return func(r chi.Router) {
		r.HandleFunc("/meal/{id}", MealFunc)
		r.HandleFunc("/meal", MealFunc)
	}
}

func MealHandler(writer http.ResponseWriter, request *http.Request) {
	repoInterface := request.Context().Value(repository.ContextKeyRepository)
	sessionRepo, ok := repoInterface.(repository.Repository)
	if !ok {
		panic("could not fetch repo for session")
	}
	MealRouter(writer, request, services.NewMealService(sessionRepo.MealRepository()), services.NewMealFoodService(sessionRepo.MealFoodRepository()))
}

func MealRouter(writer http.ResponseWriter, request *http.Request, service services.MealService, mFService services.MealFoodService) {
	routePattern := chi.RouteContext(request.Context()).RoutePattern()
	switch request.Method {
	case http.MethodGet:
		if routePattern == "/api/meal/{userId}/{date}" {
			GetUserMealsByDate(writer, request, service, mFService)
		}
	case http.MethodPost:
		if routePattern == "/api/meal" {
			CreateMeal(writer, request, service, mFService)
		}
	case http.MethodPut:
		if routePattern == "/api/meal" {
			UpdateMeal(writer, request, service, mFService)
		}
	case http.MethodDelete:
		if routePattern == "/api/meal/{id}" {
			DeleteMeal(writer, request, service, mFService)
		}
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateMeal(writer http.ResponseWriter, request *http.Request, service services.MealService, mFService services.MealFoodService) {
	var meal entity.Meal
	if err := json.NewDecoder(request.Body).Decode(&meal); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}
	createdMeal, err := service.CreateEntry(request.Context(), &meal)
	if err != nil {
		http.Error(writer, "Failed to create meal", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(createdMeal)
	if err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetUserMealsByDate(writer http.ResponseWriter, request *http.Request, service services.MealService, mFService services.MealFoodService) {
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
	if err := json.NewEncoder(writer).Encode(entries); err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func UpdateMeal(writer http.ResponseWriter, request *http.Request, service services.MealService, mFService services.MealFoodService) {
	var meal entity.Meal
	if err := json.NewDecoder(request.Body).Decode(&meal); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := service.UpdateEntry(request.Context(), &meal); err != nil {
		http.Error(writer, "Failed to update meal", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(writer).Encode(meal); err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func DeleteMeal(writer http.ResponseWriter, request *http.Request, service services.MealService, mFService services.MealFoodService) {
	id := chi.URLParam(request, "id")
	if err := service.DeleteEntry(request.Context(), id); err != nil {
		http.Error(writer, "Failed to delete meal", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
