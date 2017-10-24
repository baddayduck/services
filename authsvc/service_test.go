package authsvc_test

import (
	"context"
	"testing"

	"github.com/baddayduck/services/authsvc"
)

func TestCreateSession(t *testing.T) {
	ctx := context.Background()
	svc := authsvc.NewInmemService()
	err := svc.CreateSession(ctx, authsvc.Session{ID: "1", Username: "bdd"})
	if err != nil {
		t.Errorf("failed to create session: %s", err)
	}
}

func TestGetSession(t *testing.T) {
	ctx := context.Background()
	svc := authsvc.NewInmemService()
	svc.CreateSession(ctx, authsvc.Session{ID: "1", Username: "bdd"})

	session, err := svc.GetSession(ctx, "1")
	if err != nil {
		t.Errorf("failed to get session: %s", err)
	}

	if session.ID != "1" || session.Username != "bdd" {
		t.Errorf("session did not match")
	}
}

func TestDeleteSession(t *testing.T) {
	ctx := context.Background()
	svc := authsvc.NewInmemService()
	svc.CreateSession(ctx, authsvc.Session{ID: "1", Username: "bdd"})

	if err := svc.DeleteSession(ctx, "1"); err != nil {
		t.Errorf("failed to delete session: %s", err)
	}

	_, err := svc.GetSession(ctx, "1")
	if err != authsvc.ErrNotFound {
		t.Errorf("session was not deleted")
	}
}
