package db

import (
	"time"

	user "github.com/baddayduck/services/usersvc/proto/account"
)

var (
	Url = "root:root@tcp(127.0.0.1:3306)/user"
)

func Init() {

}

func Create(user *user.User) error {
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()
	return nil
}

func Read(id string) (*user.User, error) {
	user := &user.User{}
	// do a query
	return user, nil
}
