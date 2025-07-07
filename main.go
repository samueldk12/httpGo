package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main(){
	r := chi.NewMux()
	
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/horario", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		fmt.Fprintln(w, now)
	})

	r.Route("/api",func(r chi.Router) {
		r.Route("/v1",func(r chi.Router) {
			r.Get(
				"/users", func(w http.ResponseWriter, r *http.Request) {
				},
			)
		})

		r.Route("/v2",func(r chi.Router) {
			r.Get(
				"/users", func(w http.ResponseWriter, r *http.Request) {
				},
			)
		})

		r.With(middleware.RealIP).Get("/users", func(w http.ResponseWriter, r *http.Request) {})

		r.Group(func(r chi.Router) {
			r.Use(middleware.BasicAuth("",
			map[string]string{
				"admin": "admin",
			}))

			r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "ping")
			})
		})
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}