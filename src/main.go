package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	Id       int     `json:"id"`
	Username string  `json:"name"`
	Age      int     `json:"age"`
	Friends  []*User `json:"friends"`
}

type Storage struct {
	Users map[int]*User
}

type Friend struct {
	SourceId int `json:"source_id"`
	TargetId int `json:"target_id"`
}

type DeleteUser struct {
	TargetId int `json:"target_id"`
}

type NewAge struct {
	New_age int `json:"new_age"`
}

func main() {
	e := chi.NewRouter()
	storage := Storage{make(map[int]*User)}
	e.Use(middleware.Logger)
	e.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HELLO"))
	})
	e.Route("/users", func(e chi.Router) {
		e.Post("/create", storage.create)
		e.Get("/show", storage.get)
		e.Post("/make_friends", storage.makeFriends)
		e.Delete("/delete_user", storage.Delete)
		e.Route("/show_friends", func(e chi.Router) {
			e.Get("/{userId}", storage.showFriends)
		})
		e.Route("/set_age", func(e chi.Router) {
			e.Put("/{userId}", storage.setAge)
		})
	})
	http.ListenAndServe(":3333", e)
}

func (s Storage) create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal()
			return
		}
		var u *User
		if err := json.Unmarshal(content, &u); err != nil {
			log.Fatal(err)
			return
		}
		s.Users[u.Id] = u
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created " + u.Username))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s Storage) get(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response := ""
		for _, user := range s.Users {
			response += user.print()
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(response))
		return
	}
	writer.WriteHeader(http.StatusBadRequest)
}

func (u *User) print() string {
	return fmt.Sprintf("Id is %v, name is %s, age is %d\n", u.Id, u.Username, u.Age)
}

func (s Storage) makeFriends(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	if request.Method == "POST" {
		content, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		sender := new(User)
		receiver := new(User)
		var f *Friend
		if err := json.Unmarshal(content, &f); err != nil {
			fmt.Println(err)
			return
		}
		sender.Id = f.SourceId
		receiver.Id = f.TargetId
		if _, ok := s.Users[sender.Id]; !ok {
			if _, ok := s.Users[receiver.Id]; !ok {
				writer.WriteHeader(http.StatusNotFound)
				writer.Write([]byte("User with id:" + strconv.Itoa(receiver.Id) + " is not found!"))
				return
			}
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte("User with id:" + strconv.Itoa(sender.Id) + " is not found!"))
			return
		}
		if !checkIfIsFriend(s.Users[sender.Id].Friends, s.Users[receiver.Id]) {
			s.Users[sender.Id].Friends = append(s.Users[sender.Id].Friends, s.Users[receiver.Id])
			s.Users[receiver.Id].Friends = append(s.Users[receiver.Id].Friends, s.Users[sender.Id])
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("Id:" + strconv.Itoa(sender.Id) + " and id:" + strconv.Itoa(receiver.Id) + " are friends"))
			return
		}
		writer.Write([]byte("Id:" + strconv.Itoa(sender.Id) + " and id:" + strconv.Itoa(receiver.Id) + " already friends"))
		return
	}
	writer.WriteHeader(http.StatusBadRequest)
}

func checkIfIsFriend(friends []*User, u *User) bool {
	for _, i := range friends {
		if i == u {
			return true
		}
	}
	return false
}

func (s Storage) showFriends(w http.ResponseWriter, r *http.Request) {
	u := new(User)
	if r.Method == "GET" {
		id := chi.URLParam(r, "userId")
		realId, err := strconv.Atoi(id)
		if err != nil {
			log.Fatal(err)
			return
		}
		u.Id = realId
		response := ""
		for i := range s.Users[u.Id].Friends {
			response += fmt.Sprintf("%v Friend's id: %v,  name: %s, age: %d\n", i+1, s.Users[u.Id].Friends[i].Id, s.Users[u.Id].Friends[i].Username, s.Users[u.Id].Friends[i].Age)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s Storage) Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "DELETE" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		user := new(User)
		var target_user *DeleteUser
		if err := json.Unmarshal(content, &target_user); err != nil {
			fmt.Println(err)
			return
		}
		user.Id = target_user.TargetId
		if _, ok := s.Users[user.Id]; !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		delete(s.Users, user.Id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("USER WAS DELETED"))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s Storage) DeleteFromFriends(id1, id2 int) {
	if _, ok := s.Users[id1]; !ok {
		return
	}
	for i, _ := range s.Users[id1].Friends {
		if i == id2 {

		}
	}
}

func (s Storage) setAge(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "PUT" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		u := new(User)
		userid := chi.URLParam(r, "userId")
		realId, err := strconv.Atoi(userid)
		if err != nil {
			log.Fatal(err)
			return
		}
		u.Id = realId
		if _, ok := s.Users[u.Id]; !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var new_age *NewAge
		if err := json.Unmarshal(content, &new_age); err != nil {
			fmt.Println(err)
			return
		}
		s.Users[u.Id].Age = new_age.New_age
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("STATUS EPTA"))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
