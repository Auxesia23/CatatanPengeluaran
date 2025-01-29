package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Auxesia23/CatatanPengeluaran/database"
	"github.com/Auxesia23/CatatanPengeluaran/handlers"
	"github.com/Auxesia23/CatatanPengeluaran/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.InitDB()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userHandler := handlers.NewUserHandler(db)
	transactionHandler := handlers.NewTransactionHandler(db)
	methodHandler := handlers.NewMethodHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)

	r.Route("/user", func(r chi.Router) {
		r.Use(middlewares.JWTAuthMiddleware)
		r.Get("/me", userHandler.GetUser)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
	})

	r.Route("/transaction", func(r chi.Router) {
		r.Use(middlewares.JWTAuthMiddleware)
		r.Post("/create", transactionHandler.CreateTransaction)

		r.Get("/method", methodHandler.GetMethods)
		r.Get("/method/{id}", methodHandler.GetMethod)
		r.Post("/method/create", methodHandler.CreateMethod)
		r.Delete("/method/delete/{id}", methodHandler.DeleteMethod)
		r.Put("/method/update/{id}", methodHandler.UpdateMethod)

		r.Get("/category", categoryHandler.GetCategories)
		r.Get("/category/{id}", categoryHandler.GetCategory)
		r.Post("/category/create", categoryHandler.CreateCategory)
		r.Put("/category/update/{id}", categoryHandler.UpdateCategory)
		r.Delete("/category/delete/{id}", categoryHandler.DeleteCategory)
	})

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}
	log.Println("Server is running on port " + os.Getenv("PORT"))
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
