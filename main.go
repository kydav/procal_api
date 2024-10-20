package main

import (
	"context"
	"log"
	"net/http"
	"procal/api/services"
	"procal/api/wrappers/FatSecretWrapper"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("local.env")
	//err := godotenv.Load("procal.env")
	if err != nil {
		log.Panic(err)
	}
	fatSecretWrapper := FatSecretWrapper.NewFatSecretWrapper()
	nutritionService := services.NewNutritionService(fatSecretWrapper)
	nutritionService.FindById(context.Background())

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(":3000", r)
}
