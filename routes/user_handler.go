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
		returnError(writer, "Bad Request", http.StatusBadRequest, nil)
	}
}

func UserByIdFinder(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	id := chi.URLParam(request, "id")
	if id == "" {
		returnError(writer, "Invalid User ID", http.StatusBadRequest, nil)
		return
	}

	user, err := service.FindById(request.Context(), id)
	if err != nil {
		returnError(writer, "User Not Found", http.StatusNotFound, err)
		return
	}
	if user.ID == "" {
		returnError(writer, "User Not Found", http.StatusNotFound, nil)
		return
	}
	returnSuccess(writer, user)
}

func UserByEmailFinder(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	email := chi.URLParam(request, "email")
	if email == "" {
		returnError(writer, "Email is required", http.StatusBadRequest, nil)
		return
	}

	user, err := service.FindByEmail(request.Context(), email)
	if err != nil {
		returnError(writer, "User Not Found", http.StatusNotFound, err)
		return
	}
	returnSuccess(writer, user)
}

func UserByFirebaseUidFinder(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	firebaseUid := chi.URLParam(request, "firebaseUid")
	if firebaseUid == "" {
		returnError(writer, "Firebase UID is required", http.StatusBadRequest, nil)
		return
	}

	user, err := service.FindByFirebaseUid(request.Context(), firebaseUid)
	if err != nil {
		returnError(writer, "User Not Found", http.StatusNotFound, err)
		return
	}
	returnSuccess(writer, user)
}

func CreateUser(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	var user entity.User
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		returnError(writer, "Invalid Request Body", http.StatusBadRequest, err)
		return
	}

	createdUser, err := service.Create(request.Context(), user)
	if err != nil {
		returnError(writer, "Error Creating User", http.StatusInternalServerError, err)
		return
	}
	returnSuccess(writer, createdUser)
}

func UpdateUser(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	var user entity.User
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&user)
	if err != nil {
		returnError(writer, "Invalid Request Body", http.StatusBadRequest, err)
		return
	}

	updatedUser, err := service.Update(request.Context(), user)
	if err != nil {
		returnError(writer, "Error Updating User", http.StatusInternalServerError, err)
		return
	}
	returnSuccess(writer, updatedUser)
}

func DeleteUser(writer http.ResponseWriter, request *http.Request, service services.UserService) {
	id := chi.URLParam(request, "id")
	if id == "" {
		returnError(writer, "Invalid User ID", http.StatusBadRequest, nil)
		return
	}

	if err := service.Delete(request.Context(), id); err != nil {
		returnError(writer, "Error Deleting User", http.StatusInternalServerError, err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
