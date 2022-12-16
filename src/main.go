package main

import (
	"net/http"
	"serv/src/database"
	"serv/src/storage"
	"serv/src/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	database.CreateTable()
	e := chi.NewRouter()
	storage := storage.Storage{Users: make(map[int]*user.User)}
	e.Use(middleware.Logger)
	e.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HELLO"))
	})
	e.Route("/users", func(e chi.Router) {
		e.Post("/create", storage.Create)
		e.Get("/show", storage.Get)
		e.Post("/make_friends", storage.MakeFriends)
		e.Delete("/delete_user", storage.Delete)
		e.Route("/show_friends", func(e chi.Router) {
			e.Get("/{userId}", storage.ShowFriends)
		})
		e.Route("/set_age", func(e chi.Router) {
			e.Put("/{userId}", storage.SetAge)
		})
	})
	http.ListenAndServe(":3333", e)
}
