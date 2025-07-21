package routes

import (
	"encoding/json"
	"net/http"
	"procal/entity"
	"procal/repository"
	"procal/services"
	"strconv"

	"github.com/go-chi/chi"
)

func UserRoutes() func(chi.Router) {
	UserFunc := func(writer http.ResponseWriter, request *http.Request) {
		UserHandler(writer, request)
	}
	return func(r chi.Router) {
		r.HandleFunc("/user/{id}", UserFunc)
		r.HandleFunc("/user/email/{email}", UserFunc)
		r.HandleFunc("/user/create", UserFunc)
		r.HandleFunc("/user/update", UserFunc)
		r.HandleFunc("/user/delete/{id}", UserFunc)
	}
}

func UserHandler(writer http.ResponseWriter, request *http.Request) {
	repoInterface := request.Context().Value("Repository")
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
	case http.MethodPost:
		if routePattern == "/api/user/create" {
			CreateUser(writer, request, service)
			return
		}
	case http.MethodPut:
		if routePattern == "/api/user/update" {
			UpdateUser(writer, request, service)
			return
		}
	case http.MethodDelete:
		if routePattern == "/api/user/delete/{id}" {
			DeleteUser(writer, request, service)
			return
		}
	default:
		http.Error(writer, "Bad Request", http.StatusBadRequest)
	}
}

func UserByIdFinder(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	id := chi.URLParam(request, "id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(writer, "Invalid User ID", http.StatusBadRequest)
		return
	}

	user, err := service.FindById(userId)
	if err != nil {
		http.Error(writer, "User Not Found", http.StatusNotFound)
		return
	}
	if user.ID == "" {
		http.Error(writer, "User Not Found", http.StatusNotFound)
		return
	}
	json.NewEncoder(writer).Encode(user)

}
func UserByEmailFinder(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	email := chi.URLParam(request, "email")
	if email == "" {
		http.Error(writer, "Email is required", http.StatusBadRequest)
		return
	}

	user, err := service.FindByEmail(email)
	if err != nil {
		http.Error(writer, "User Not Found", http.StatusNotFound)
		return
	}

	json.NewEncoder(writer).Encode(user)
}
func CreateUser(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	var user entity.User
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		http.Error(writer, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	createdUser, err := service.Create(user)
	if err != nil {
		http.Error(writer, "Error Creating User", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(createdUser)
}

func UpdateUser(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	var user entity.User
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		http.Error(writer, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	updatedUser, err := service.Update(user)
	if err != nil {
		http.Error(writer, "Error Updating User", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(updatedUser)
}
func DeleteUser(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	id := chi.URLParam(request, "id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(writer, "Invalid User ID", http.StatusBadRequest)
		return
	}

	if err := service.Delete(userId); err != nil {
		http.Error(writer, "Error Deleting User", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
