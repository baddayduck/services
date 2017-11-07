package handler

import (
	"context"
	"strings"
	"time"

	"github.com/baddayduck/services/authsvc/proto/auth"
)

type Auth struct{}

func random(i int) string {
	return "randomstring"
}

func (s *Auth) Login(ctx context.Context, req *auth.LoginRequest, rsp *auth.LoginResponse) error {
	username := strings.ToLower(req.Username)
	// email := strings.ToLower(req.Email)

	// TODO salt/hash compute

	sess := &auth.Session{
		Id:       random(128),
		Username: username,
		Created:  time.Now().Unix(),
		Expires:  time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	// TODO save the session
	rsp.Session = sess
	return nil
}

func (s *Auth) Logout(ctx context.Context, req *auth.LogoutRequest, rsp *auth.LogoutResponse) error {
	// TODO delete the session
	return nil
}

func (s *Auth) ReadSession(ctx context.Context, req *auth.ReadSessionRequest, rsp *auth.ReadSessionResponse) error {
	// TODO query the session
	rsp.Session = &auth.Session{}
	return nil
}
