package storage

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
	"serv/src/delete_user"
	"serv/src/make_friends"
	"serv/src/new_age"
	"serv/src/user"
	"strconv"
)

type Storage struct {
	Users map[int]*user.User
}

func (s Storage) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal()
			return
		}
		var u *user.User
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

func (s Storage) Get(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response := ""
		for _, u := range s.Users {
			response += u.Print()
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(response))
		return
	}
	writer.WriteHeader(http.StatusBadRequest)
}

func (s Storage) MakeFriends(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	if request.Method == "POST" {
		content, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		sender := new(user.User)
		receiver := new(user.User)
		var f *make_friends.Friend
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

func checkIfIsFriend(friends []*user.User, u *user.User) bool {
	for _, i := range friends {
		if i == u {
			return true
		}
	}
	return false
}

func (s Storage) ShowFriends(w http.ResponseWriter, r *http.Request) {
	u := new(user.User)
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
		user := new(user.User)
		var targetUser *delete_user.DeleteUser
		if err := json.Unmarshal(content, &targetUser); err != nil {
			fmt.Println(err)
			return
		}
		user.Id = targetUser.TargetId
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

func (s Storage) SetAge(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "PUT" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		u := new(user.User)
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
		var newAge *new_age.NewAge
		if err := json.Unmarshal(content, &newAge); err != nil {
			fmt.Println(err)
			return
		}
		s.Users[u.Id].Age = newAge.NewAge
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("STATUS EPTA"))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
