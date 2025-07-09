package main

import (
	"log"
	"net/http"
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

	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Group(routes.NutritionRoutes())
	})
	http.ListenAndServe("0.0.0.0:8000", r)

}
