package handler

import (
	"context"

	"github.com/baddayduck/services/usersvc/db"
	"github.com/baddayduck/services/usersvc/proto/account"
)

type Account struct{}

func (s *Account) Create(ctx context.Context, req *account.CreateRequest, rsp *account.CreateResponse) error {
	return db.Create(req.User)
}

func (s *Account) Read(ctx context.Context, req *account.ReadRequest, rsp *account.ReadResponse) error {
	user, err := db.Read(req.Id)
	if err != nil {
		return err
	}
	rsp.User = user
	return nil
}
