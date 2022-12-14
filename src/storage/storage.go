package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"serv/src/database"
	"serv/src/delete_user"
	"serv/src/make_friends"
	"serv/src/new_age"
	"serv/src/user"
	"strconv"

	
	"github.com/go-chi/chi/v5"
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
		if !database.CheckUser(u.Id){
			if err = database.AddUser(content); err != nil{
				log.Fatal(err)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Created " + u.Username))
			return
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s Storage) Get(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response := ""
		for _, u := range database.GetUsers(){
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
		var f *make_friends.Friend
		if err := json.Unmarshal(content, &f); err != nil {
			fmt.Println(err)
			return
		}
		sender := database.GetUser(f.SourceId)
		receiver := database.GetUser(f.TargetId)
		if !database.CheckUser(sender.Id){
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte("User with id:" + strconv.Itoa(sender.Id) + " is not found!"))
			return
		}
		if !database.CheckUser(receiver.Id){
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte("User with id:" + strconv.Itoa(receiver.Id) + " is not found!"))
			return
		}
		if !database.CheckIfIsFriend(sender.Id, receiver.Id){
			if err = database.AddFriends(sender.Id, receiver.Id); err != nil{
				log.Fatal(err)
			}
			if err = database.AddFriends(receiver.Id, sender.Id); err != nil{
				log.Fatal(err)
			}
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("Id:" + strconv.Itoa(sender.Id) + " and id:" + strconv.Itoa(receiver.Id) + " are friends"))
			return
		}
		writer.WriteHeader(http.StatusAlreadyReported)
		writer.Write([]byte("Id:" + strconv.Itoa(sender.Id) + " and id:" + strconv.Itoa(receiver.Id) + " are already friends!"))
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
		friends, err := database.GetFriends(u.Id)
		if err != nil{
			log.Fatal(err)
			return
		}
		response := ""
		for _, val := range friends{
			fmt.Println(val)
			response += database.GetUser(val).Print()
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
		if !database.CheckUser(user.Id){
			w.WriteHeader(http.StatusNotFound)
			return
		}
		friends, err := database.GetFriends(user.Id)
		s.DeleteFromFriends(user.Id, friends)
		if err = database.DeleteUser(user.Id); err != nil{
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User was succesfully deleted!"))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func remove(a[]*user.User, ind int) []*user.User{
	new := make([]*user.User, 0)
	new = append(new, a[:ind]...)
	return append(new, a[ind+1:]...)
}

func (s Storage) DeleteFromFriends(id1 int, arr []int){
	if !database.CheckUser(id1){
		return
	}
	for i := range arr{
		database.DeleteFromFriends(arr[i], id1)
	}
}


func (s Storage) GetAllFriendsId(u *user.User)[]int{
	if _, ok := s.Users[u.Id];!ok{
		return nil
	}
	friends, err := database.GetFriends(u.Id)
	if err!=nil{
		log.Fatal(err)
		return nil
	}
	return friends
}

func (s Storage) SetAge(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "PUT" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		userid := chi.URLParam(r, "userId")
		realId, err := strconv.Atoi(userid)
		if err != nil {
			log.Fatal(err)
			return
		}
		if !database.CheckUser(realId){
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var newAge *new_age.NewAge
		if err := json.Unmarshal(content, &newAge); err != nil {
			fmt.Println(err)
			return
		}
		database.SetAge(realId, newAge.NewAge)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("New age is set"))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
