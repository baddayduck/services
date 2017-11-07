package db

import (
	"time"

	auth "github.com/baddayduck/services/authsvc/proto/auth"
)

var (
	Url = "root:root@tcp(127.0.0.1:3306)/auth"
)

func Init() {

}

func ReadSession(id string) (*auth.Session, error) {
	sess := &auth.Session{}

	// TODO query for the session

	return sess, nil
}

func CreateSession(sess *auth.Session) error {
	if sess.Created == 0 {
		sess.Created = time.Now().Unix()
	}

	if sess.Expires == 0 {
		sess.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	// TODO persist the session
	return nil
}

func DeleteSession(id string) error {
	// TODO delete the session
	return nil
}
