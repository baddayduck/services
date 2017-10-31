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

func TestCreateDuplicateSession(t *testing.T) {
	ctx := context.Background()
	svc := authsvc.NewInmemService()
	svc.CreateSession(ctx, authsvc.Session{ID: "1", Username: "bdd"})
	err := svc.CreateSession(ctx, authsvc.Session{ID: "1", Username: "bdd"})
	if err != authsvc.ErrAlreadyExists {
		t.Errorf("should not allow duplicate sessions")
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

func TestGetNonExistentSession(t *testing.T) {
	ctx := context.Background()
	svc := authsvc.NewInmemService()

	_, err := svc.GetSession(ctx, "1")
	if err != authsvc.ErrNotFound {
		t.Errorf("did not get ErrNotFound")
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

func TestDeleteNonExistentSession(t *testing.T) {
	ctx := context.Background()
	svc := authsvc.NewInmemService()

	if err := svc.DeleteSession(ctx, "1"); err != authsvc.ErrNotFound {
		t.Errorf("should not be able to delete a non-existent session")
	}
}
