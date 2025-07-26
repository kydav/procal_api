package routes

import (
	"encoding/json"
	"net/http"
	"procal/entity"
	"procal/repository"
	"procal/services"

	"github.com/go-chi/chi"
)

func GoalRoutes() func(chi.Router) {
	GoalFunc := func(writer http.ResponseWriter, request *http.Request) {
		GoalHandler(writer, request)
	}
	return func(r chi.Router) {
		r.HandleFunc("/goal", GoalFunc)
		r.HandleFunc("/goal/{id}", GoalFunc)
	}
}
func GoalHandler(writer http.ResponseWriter, request *http.Request) {
	repoInterface := request.Context().Value(repository.ContextKeyRepository)
	sessionRepo, ok := repoInterface.(repository.Repository)
	if !ok {
		http.Error(writer, "could not fetch repo for session", http.StatusInternalServerError)
		return
	}
	goalService := services.NewGoalService(sessionRepo.GoalRepository())
	GoalRouter(writer, request, goalService)
}
func GoalRouter(writer http.ResponseWriter, request *http.Request, service services.GoalService) {
	routePattern := chi.RouteContext(request.Context()).RoutePattern()
	switch request.Method {
	case http.MethodPost:
		if routePattern == "/api/goal" {
			CreateGoal(writer, request, service)
			return
		}
	case http.MethodGet:
		if routePattern == "/api/goal/{id}" {
			GetGoalByID(writer, request, service)
			return
		}
	case http.MethodPut:
		if routePattern == "/api/goal" {
			UpdateGoal(writer, request, service)
			return
		}
	case http.MethodDelete:
		if routePattern == "/api/goal/{id}" {
			DeleteGoal(writer, request, service)
			return
		}
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func CreateGoal(writer http.ResponseWriter, request *http.Request, service services.GoalService) {
	var goal entity.Goal
	if err := json.NewDecoder(request.Body).Decode(&goal); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := service.CreateGoal(request.Context(), &goal); err != nil {
		http.Error(writer, "Failed to create goal", http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(goal)
}
func GetGoalByID(writer http.ResponseWriter, request *http.Request, service services.GoalService) {
	id := chi.URLParam(request, "id")
	goal, err := service.GetGoalByID(request.Context(), id)
	if err != nil {
		http.Error(writer, "Goal not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(writer).Encode(goal)
}
func UpdateGoal(writer http.ResponseWriter, request *http.Request, service services.GoalService) {
	var goal entity.Goal
	if err := json.NewDecoder(request.Body).Decode(&goal); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := service.UpdateGoal(request.Context(), &goal); err != nil {
		http.Error(writer, "Failed to update goal", http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(goal)
}
func DeleteGoal(writer http.ResponseWriter, request *http.Request, service services.GoalService) {
	id := chi.URLParam(request, "id")
	if err := service.DeleteGoal(request.Context(), id); err != nil {
		http.Error(writer, "Failed to delete goal", http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}
