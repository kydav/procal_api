package routes

import (
	"encoding/json"
	"net/http"

	"procal/entity"
	"procal/repository"
	"procal/services"

	"github.com/go-chi/chi"
)

func UserRoutes() func(chi.Router) {
	UserFunc := func(writer http.ResponseWriter, request *http.Request) {
		UserHandler(writer, request)
	}
	return func(r chi.Router) {
		r.HandleFunc("/user/{id}", UserFunc)
		r.HandleFunc("/user/email/{email}", UserFunc)
		r.HandleFunc("/user/uid/{firebaseUid}", UserFunc)
		r.HandleFunc("/user", UserFunc)
		r.HandleFunc("/user/{id}", UserFunc)
	}
}

func UserHandler(writer http.ResponseWriter, request *http.Request) {
	repoInterface := request.Context().Value(repository.ContextKeyRepository)
	sessionRepo, ok := repoInterface.(repository.Repository)
	if !ok {
		panic("could not fetch repo for session")
	}
	UserRouter(writer, request, services.NewUserService(sessionRepo.UserRepository()))
}

func UserRouter(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	routePattern := chi.RouteContext(request.Context()).RoutePattern()
	switch request.Method {
	case http.MethodGet:
		if routePattern == "/api/user/{id}" {
			UserByIdFinder(writer, request, service)
			return
		}
		if routePattern == "/api/user/email/{email}" {
			UserByEmailFinder(writer, request, service)
			return
		}
		if routePattern == "/api/user/uid/{firebaseUid}" {
			UserByFirebaseUidFinder(writer, request, service)
			return
		}
	case http.MethodPost:
		if routePattern == "/api/user" {
			CreateUser(writer, request, service)
			return
		}
	case http.MethodPut:
		if routePattern == "/api/user" {
			UpdateUser(writer, request, service)
			return
		}
	case http.MethodDelete:
		if routePattern == "/api/user/{id}" {
			DeleteUser(writer, request, service)
			return
		}
	default:
		http.Error(writer, "Bad Request", http.StatusBadRequest)
	}
}

func UserByIdFinder(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	id := chi.URLParam(request, "id")
	if id == "" {
		http.Error(writer, "Invalid User ID", http.StatusBadRequest)
		return
	}

	user, err := service.FindById(request.Context(), id)
	if err != nil {
		http.Error(writer, "User Not Found", http.StatusNotFound)
		return
	}
	if user.ID == "" {
		http.Error(writer, "User Not Found", http.StatusNotFound)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(user)
}

func UserByEmailFinder(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	email := chi.URLParam(request, "email")
	if email == "" {
		http.Error(writer, "Email is required", http.StatusBadRequest)
		return
	}

	user, err := service.FindByEmail(request.Context(), email)
	if err != nil {
		http.Error(writer, "User Not Found", http.StatusNotFound)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	json.NewEncoder(writer).Encode(user)
}

func UserByFirebaseUidFinder(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	firebaseUid := chi.URLParam(request, "firebaseUid")
	if firebaseUid == "" {
		http.Error(writer, "Firebase UID is required", http.StatusBadRequest)
		return
	}

	user, err := service.FindByFirebaseUid(request.Context(), firebaseUid)
	if err != nil {
		http.Error(writer, "User Not Found", http.StatusNotFound)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	json.NewEncoder(writer).Encode(user)
}

func CreateUser(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	var user entity.User
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		http.Error(writer, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	createdUser, err := service.Create(request.Context(), user)
	if err != nil {
		http.Error(writer, "Error Creating User", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(createdUser)
}

func UpdateUser(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	var user entity.User
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(writer, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	updatedUser, err := service.Update(request.Context(), user)
	if err != nil {
		http.Error(writer, "Error Updating User", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	json.NewEncoder(writer).Encode(updatedUser)
}

func DeleteUser(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	id := chi.URLParam(request, "id")
	if id == "" {
		http.Error(writer, "Invalid User ID", http.StatusBadRequest)
		return
	}

	if err := service.Delete(request.Context(), id); err != nil {
		http.Error(writer, "Error Deleting User", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
