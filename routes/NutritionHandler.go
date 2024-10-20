package routes

import (
	"net/http"
)

type NutritionHandler struct {
}

func (f NutritionHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	// err := json.NewEncoder(w).Encode(listBooks())
	// if err != nil {
	// 	http.Error(w, "Internal error", http.StatusInternalServerError)
	// 	return
	// }
}
