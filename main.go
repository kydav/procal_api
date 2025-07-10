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
	// fat_secret_client_id := os.Getenv("FAT_SECRET_CLIENT_ID")
	// if fat_secret_client_id == "" {
	// 	log.Panic("FAT_SECRET_CLIENT_ID environment variable is not set")
	// }
	// fat_secret_client_secret := os.Getenv("FAT_SECRET_CLIENT_SECRET")
	// if fat_secret_client_secret == "" {
	// 	log.Panic("FAT_SECRET_CLIENT_SECRET environment variable is not set")
	// }
	// // Initialize the FatSecretWrapper with the client ID and secret
	// fatSecretWrapper := routes.NewFatSecretWrapper(fat_secret_client_id, fat_secret_client_secret)

	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Group(routes.NutritionRoutes())
	})
	http.ListenAndServe("0.0.0.0:8000", r)

}
