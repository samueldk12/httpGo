package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Response struct {
	Error string `json:"error,omitempty"`
	Data any `json:"data"`
}

func sendJSON(w http.ResponseWriter, resp Response, status int){
	data, err := json.Marshal(resp)

	if err != nil {
		fmt.Println("error ao fazer um marshal de json:", err)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		fmt.Println("error ao enviar a resposta")
		return
	}
}


type User struct{
	Username string 
	Password string `json:"-"`
	Role string
	ID int64 `json:"id,string"`
}

func main(){
	r := chi.NewMux()
	
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	db := map[int64]User{
		1: {
			Username: "admin",
			Password: "admin",
			Role: "admin",
			ID: 1,
		},
	}
	
	r.Group(func(r chi.Router) {
		r.Use(jsonMiddleware)
		r.Get("/users/{id:[0-9]+}", handleGetUsers(db))
		r.Post("/users", handlePostUsers(db))
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}

}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		
		next.ServeHTTP(w, r)
	})
}

func handleGetUsers(db map[int64]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		idStr := chi.URLParam(r, "id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		user, ok := db[id]

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "usuario nao encontrado"}`))
			return
		}
	
		data, err := json.Marshal(user)
		
		if err != nil {
			sendJSON(
				w,
				Response{Error: "somethin went wrong"},
				http.StatusInternalServerError,
			)
			return
		}
		
		_, _ = w.Write(data)
	
	}
}

func handlePostUsers(db map[int64]User) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		http.MaxBytesReader(w, r.Body, 10000)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr){
				sendJSON(
					w, 
					Response{Error: "body too large"},
					http.StatusRequestEntityTooLarge,
				)
				return
			}
		}

		var user User

		if err := json.Unmarshal(data, &user); err != nil {
			sendJSON(
				w, 
				Response{Error:"invalid Body"}, 
				http.StatusUnprocessableEntity,
			)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}