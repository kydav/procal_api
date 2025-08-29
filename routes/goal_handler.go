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
		returnError(writer, "could not fetch repo for session", http.StatusInternalServerError, nil)
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
		returnError(writer, "Method not allowed", http.StatusMethodNotAllowed, nil)
	}
}

func CreateGoal(writer http.ResponseWriter, request *http.Request, service services.GoalService) {
	var goal entity.Goal
	if err := json.NewDecoder(request.Body).Decode(&goal); err != nil {
		returnError(writer, "Invalid request body", http.StatusBadRequest, err)
		return
	}
	if err := service.CreateGoal(request.Context(), &goal); err != nil {
		returnError(writer, "Failed to create goal", http.StatusInternalServerError, err)
		return
	}
	returnSuccess(writer, goal)
}

func GetGoalByID(writer http.ResponseWriter, request *http.Request, service services.GoalService) {
	id := chi.URLParam(request, "id")
	goal, err := service.GetGoalByID(request.Context(), id)
	if err != nil {
		returnError(writer, "Goal not found", http.StatusNotFound, err)
		return
	}
	returnSuccess(writer, goal)
}

func UpdateGoal(writer http.ResponseWriter, request *http.Request, service services.GoalService) {
	var goal entity.Goal
	if err := json.NewDecoder(request.Body).Decode(&goal); err != nil {
		returnError(writer, "Invalid request body", http.StatusBadRequest, err)
		return
	}
	if err := service.UpdateGoal(request.Context(), &goal); err != nil {
		returnError(writer, "Failed to update goal", http.StatusInternalServerError, err)
		return
	}
	returnSuccess(writer, goal)
}

func DeleteGoal(writer http.ResponseWriter, request *http.Request, service services.GoalService) {
	id := chi.URLParam(request, "id")
	if err := service.DeleteGoal(request.Context(), id); err != nil {
		returnError(writer, "Failed to delete goal", http.StatusInternalServerError, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}
