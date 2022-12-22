package database

import (
	"database/sql"
	"fmt"
	"log"
	"serv/src/user"

	_ "github.com/lib/pq"
)

const(
	host = "localhost"
	port = 5432
	owner = "postgres"
	password = "lolkek12"
	name = "Users"
)

const(
	table_name = "users_table"
)

func Connection() string{
	return fmt.Sprintf("host = %s port = %d user = %s password = %s dbname = %s sslmode=disable", host, port, owner, password, name)
}

func CreateTable(){
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, user_data JSONB)`, table_name)
	db, err := sql.Open("postgres", Connection())
	checkErr(err)
	_, err = db.Exec(query)
	checkErr(err)
	defer db.Close()
}

func AddUser(jsonUser []byte) error{
	query := fmt.Sprintf(`INSERT INTO %s (user_data) VALUES ($1)`, table_name)
	db, err := sql.Open("postgres", Connection())
	if checkErr(err){
		return err
	}
	_, err = db.Exec(query, jsonUser)
	if checkErr(err){
		return err
	} 
	defer db.Close()
	return nil
}

func checkErr(err error) bool{
	if err != nil{
		log.Fatal(err)
		return true
	}
	return false
}

func CheckUser(id int) bool{
	db, err := sql.Open("postgres", Connection())
	checkErr(err)
	query, err := db.Query(fmt.Sprintf(`SELECT EXISTS(SELECT * FROM %s WHERE user_data->>'id' = $1)`, table_name), id)
	checkErr(err)
	var exists bool
	for	query.Next(){
		err =	query.Scan(&exists)
		checkErr(err)
	}
	fmt.Println(exists)
	return exists
}

func GetUser(id int) *user.User{
	db, err := sql.Open("postgres", Connection())
	checkErr(err)
	query, err := db.Query(fmt.Sprintf(`SELECT user_data->>'id', user_data->>'name',user_data->>'age' FROM %s WHERE user_data->>'id' = $1`, table_name), id)
	checkErr(err)
	u := new(user.User)
	for	query.Next(){
		err =	query.Scan(&u.Id, &u.Username, &u.Age)
		checkErr(err)
	}
	return u
}

func GetUsers() []*user.User{
	db, err := sql.Open("postgres", Connection())
	checkErr(err)
	query, err := db.Query(fmt.Sprintf(`SELECT user_data->>'id', user_data->>'name',user_data->>'age' FROM %s`, table_name))
	checkErr(err)
	users := make([]*user.User, 0)
	for	query.Next(){
		u := new(user.User)
		err =	query.Scan(&u.Id, &u.Username, &u.Age)
		users = append(users, u)
	}
	return users
}

func DeleteUser(id int) error{
	db, err := sql.Open("postgres", Connection())
	if checkErr(err){
		return err
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_data->>'id' = $1`, table_name)
	_, err = db.Exec(query, id)
	if checkErr(err){
		return err
	}
	fmt.Println("DELETED")
	defer db.Close()
	return nil
}

func AddFriends(senderId, receiverId int) error{
	db, err := sql.Open("postgres", Connection())
	if checkErr(err){
		return err
	}
	query := fmt.Sprintf(`update %s set user_data = jsonb_insert(user_data, '{friends,-1}', '{"id":%v}', true) where user_data->>'id'= $1`, table_name, receiverId)
	_, err = db.Exec(query,  senderId)
	if checkErr(err){
		return err
	}
	fmt.Println("DELETED")
	defer db.Close()
	return nil
}

func CheckIfIsFriend(userId, friendId int) bool{
	db, err := sql.Open("postgres", Connection())
	checkErr(err)
	query, err := db.Query(fmt.Sprintf(`select exists(select * from %s where user_data->>'id' =$1 and user_data->'friends' @> '[{"id":%v}]')`, table_name, userId), friendId)
	checkErr(err)
	var exists bool
	for	query.Next(){
		err =	query.Scan(&exists)
		checkErr(err)
	}
	return exists
}

func deleteScript(){
	db, err := sql.Open("postgres", Connection())
	checkErr(err)
	query := `create or replace function jsonb_remove_array_element(arr jsonb, element jsonb)
	returns jsonb language sql immutable as $$
		select arr- (
			select ordinality- 1
			from jsonb_array_elements(arr) with ordinality
			where value = element)::int
	$$;`
	_, err = db.Exec(query)
	defer db.Close()
}

func DeleteFromFriends(userId, friendsId int){
	deleteScript()
	db, err := sql.Open("postgres", Connection())
	checkErr(err)
	query := fmt.Sprintf(`
	update %s
	set user_data = jsonb_set(user_data, '{friends}', jsonb_remove_array_element(user_data->'friends', '{"id":%v}'))
	where user_data->'friends'  @> '[{"id":%v}]' and user_data->>'id'=$1
	returning *;`, table_name, friendsId, friendsId)
	_, err = db.Exec(query, userId)
	checkErr(err)
	defer db.Close()
}

func GetFriends(userId int) ([]int, error){
	db, err := sql.Open("postgres", Connection())
	checkErr(err)
	friends := []int{}
	query, err := db.Query(fmt.Sprintf(`select user_data->>'id' from %s where user_data->'friends' @> '[{"id":%v}]'`, table_name, userId))
	checkErr(err)
	for query.Next(){
		var i int
		if err := query.Scan(&i); err != nil{
			return nil, err
		}
		friends = append(friends, i)
	}
	return friends, nil
}
