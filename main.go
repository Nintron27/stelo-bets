package main

import (
	"chi-learning/internal/api"
	"chi-learning/internal/env"
	"log"
)

func main() {
	if err := env.Initialize(); err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(api.Start())

	// r := chi.NewRouter()

	// r.Use(cors.Handler(cors.Options{
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	// 	AllowCredentials: true,
	// 	MaxAge:           300,
	// }))

	// r.Use(middleware.Logger)
	// r.Use(middleware.URLFormat)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusCreated)
	// 	w.Write([]byte("Hello world!"))
	// })

	// r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl, err := template.ParseFiles("templates/index.html")
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}

	// 	data := struct {
	// 		Title  string
	// 		Header string
	// 	}{
	// 		Title:  "yaya | Chi",
	// 		Header: "You know it",
	// 	}

	// 	if err := tmpl.Execute(w, data); err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// })

	// http.ListenAndServe(":8080", r)
}
