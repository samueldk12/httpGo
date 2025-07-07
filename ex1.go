package main

import (
	"fmt"
	"net/http"
	"time"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		begin := time.Now()
		next.ServeHTTP(w, r)
		fmt.Fprintln(w, r.Method, time.Since(begin))
	})
}

func ex1() {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/api/users/{id}",
		func(w http.ResponseWriter, r *http.Request){
			id := r.PathValue("id")
			fmt.Fprintln(w, id)
			fmt.Fprintln(w, "Hello World!")
		},
	)


	srv := &http.Server{
		Addr:                         ":8080",
		Handler:                      Log(mux),
		DisableGeneralOptionsHandler: false,
		ReadTimeout:                  10 * time.Second,
		WriteTimeout:                 30 * time.Second,
		IdleTimeout:                  1 * time.Minute,
	}


	srv.ListenAndServe()
}