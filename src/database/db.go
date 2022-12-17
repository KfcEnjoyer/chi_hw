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
	querry := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, user_data JSONB)`, table_name)
	db, err := sql.Open("postgres", Connection())
	checkErr(err)
	_, err = db.Exec(querry)
	checkErr(err)
	defer db.Close()
}

func AddUser(jsonUser []byte) error{
	querry := fmt.Sprintf(`INSERT INTO %s (user_data) VALUES ($1)`, table_name)
	db, err := sql.Open("postgres", Connection())
	if checkErr(err){
		return err
	}
	_, err = db.Exec(querry, jsonUser)
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
	querry, err := db.Query(fmt.Sprintf(`SELECT EXISTS(SELECT * FROM %s WHERE user_data->>'id' = $1)`, table_name), id)
	checkErr(err)
	var exists bool
	for querry.Next(){
		err = querry.Scan(&exists)
		checkErr(err)
	}
	fmt.Println(exists)
	return exists
}

func GetUser(id int) *user.User{
	db, err := sql.Open("postgres", Connection())
	checkErr(err)
	querry, err := db.Query(fmt.Sprintf(`SELECT user_data->>'id', user_data->>'name',user_data->>'age' FROM %s WHERE user_data->>'id' = $1`, table_name), id)
	checkErr(err)
	u := new(user.User)
	for querry.Next(){
		err = querry.Scan(&u.Id, &u.Username, &u.Age)
		checkErr(err)
	}
	return u
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

func AddFriends(senderId int, userData []byte) error{
	db, err := sql.Open("postgres", Connection())
	if checkErr(err){
		return err
	}
	query := fmt.Sprintf(`update %s set user_data = jsonb_insert(user_data, '{friends,-1}', $1, true) where user_data->>'id'= $2`, table_name)
	_, err = db.Exec(query, userData, senderId)
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
	querry, err := db.Query(fmt.Sprintf(`select exists(select * from %s where user_data->>'id' =$1 and user_data->'friends' @> '[{"id":%v}]')`, table_name, userId), friendId)
	checkErr(err)
	var exists bool
	for querry.Next(){
		err = querry.Scan(&exists)
		checkErr(err)
	}
	return exists
}
