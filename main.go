package main

import (
	"log"
	"net/http"

	"procal/repository"
	"procal/routes"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("procal.env")
	if err != nil {
		log.Panic(err)
	}
	baseRepository := repository.NewRepository()
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(repository.BuildRepositoryWithContextMiddlware(baseRepository))
		r.Group(routes.NutritionRoutes())
		r.Group(routes.UserRoutes())
		r.Group(routes.GoalRoutes())
		r.Group(routes.MealRoutes())
	})
	err = http.ListenAndServe("0.0.0.0:8000", r)
	if err != nil {
		log.Panic(err)
	}
}
