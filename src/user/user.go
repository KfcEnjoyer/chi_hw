package user

import "fmt"

type User struct {
	Id       int     `json:"id"`
	Username string  `json:"name"`
	Age      int     `json:"age"`
	Friends  []*User `json:"friends"`
}

func (u *User) Print() string {
	return fmt.Sprintf("Id is %v, name is %s, age is %d\n", u.Id, u.Username, u.Age)
}
