package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"procal/services"
	"procal/wrappers/FatSecretWrapper"

	"github.com/go-chi/chi"
)

func NutritionRoutes() func(chi.Router) {
	NutritionFunc := func(writer http.ResponseWriter, request *http.Request) {
		NutritionHandler(writer, request)
	}
	return func(r chi.Router) {
		r.HandleFunc("/food/{id}", NutritionFunc)
		r.HandleFunc("/food/{barcode}/barcode", NutritionFunc)
		r.HandleFunc("/food/{searchQuery}/search/{page}", NutritionFunc)
	}
}

// Handles routes that will be called from the patient application
func NutritionHandler(writer http.ResponseWriter, request *http.Request) {
	fatSecretWrapper := FatSecretWrapper.NewFatSecretWrapper()
	NutritionRouter(writer, request, services.NewNutritionService(fatSecretWrapper))
}

func NutritionRouter(
	writer http.ResponseWriter,
	request *http.Request,
	service services.NutritionService,
) {
	routePattern := chi.RouteContext(request.Context()).RoutePattern()
	switch request.Method {
	case http.MethodGet:
		if routePattern == "/api/food/{id}" {
			FoodByIdFinder(writer, request, service)
			return
		}
		if routePattern == "/api/food/{barcode}/barcode" {
			FoodByBarcodeFinder(writer, request, service)
			return
		}
		if routePattern == "/api/food/{searchQuery}/search/{page}" {
			FoodFinder(writer, request, service)
			return
		}
	default:
		returnError(writer, "Bad Request", http.StatusBadRequest, errors.New("unexpected http verb"))
	}
}

func FoodByIdFinder(writer http.ResponseWriter, request *http.Request, service services.NutritionService) {
	id := chi.URLParam(request, "id")
	foodId, err := strconv.Atoi(id)
	if err != nil {
		returnError(writer, "error parsing id", http.StatusInternalServerError, err)
	}
	food, err := service.FindById(request.Context(), foodId)
	if err != nil {
		returnError(writer, "error finding food", http.StatusInternalServerError, err)
	}
	if food.Food.FoodName == "" {
		returnError(writer, "no food found", http.StatusNotFound, nil)
	} else {
		returnSuccess(writer, food)
	}
}

func FoodByBarcodeFinder(writer http.ResponseWriter, request *http.Request, service services.NutritionService) {
	barcode := chi.URLParam(request, "barcode")
	if barcode == "" {
		returnError(writer, "error parsing barcode", http.StatusInternalServerError, nil)
	}
	food, err := service.FindByBarcode(request.Context(), barcode)
	if err != nil {
		returnError(writer, "error finding food", http.StatusInternalServerError, err)
	}
	if food.Food.FoodName == "" {
		returnError(writer, "no food found", http.StatusNotFound, nil)
	} else {
		returnSuccess(writer, food)
	}
}

func FoodFinder(writer http.ResponseWriter, request *http.Request, service services.NutritionService) {
	searchQuery := chi.URLParam(request, "searchQuery")
	if searchQuery == "" {
		returnError(writer, "error parsing searchQuery", http.StatusInternalServerError, nil)
	}
	searchQuerySanitized := strings.ReplaceAll(searchQuery, "\u2019", "'")
	page := chi.URLParam(request, "page")
	foods, err := service.SearchByFoodName(request.Context(), searchQuerySanitized, page)
	if err != nil {
		returnError(writer, "error finding food", http.StatusInternalServerError, err)
	}
	if len(foods.TotalResults) == 0 {
		returnError(writer, "no foods found", http.StatusNotFound, nil)
	} else {
		returnSuccess(writer, foods)
	}
}

func returnSuccess(writer http.ResponseWriter, responseData interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(responseData); err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func returnError(writer http.ResponseWriter, errorMessage string, httpStatus int, err error) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpStatus)
	var resp ErrorResponse

	if err != nil {
		resp.Error = err.Error()
		resp.Message = errorMessage
	} else {
		resp.Error = errorMessage
	}
	jsonResp, _ := json.Marshal(resp)
	_, err = writer.Write(jsonResp)
	if err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
	}
}
