package usersvc_test

import (
	"context"
	"testing"

	"github.com/baddayduck/services/usersvc"
)

func TestAddUser(t *testing.T) {
	ctx := context.Background()
	svc := usersvc.NewInmemService()
	err := svc.AddUser(ctx, usersvc.User{ID: "1", Name: "bdd"})
	if err != nil {
		t.Errorf("failed to create session: %s", err)
	}
}

func TestCreateDuplicateUser(t *testing.T) {
	ctx := context.Background()
	svc := usersvc.NewInmemService()
	svc.AddUser(ctx, usersvc.User{ID: "1", Name: "bdd"})
	err := svc.AddUser(ctx, usersvc.User{ID: "1", Name: "bdd"})
	if err != usersvc.ErrAlreadyExists {
		t.Errorf("should not allow duplicate users")
	}
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	svc := usersvc.NewInmemService()
	svc.AddUser(ctx, usersvc.User{ID: "1", Name: "bdd"})

	user, err := svc.GetUser(ctx, "1")
	if err != nil {
		t.Errorf("failed to get user: %s", err)
	}

	if user.ID != "1" || user.Name != "bdd" {
		t.Errorf("session did not match")
	}
}

func TestGetNonExistentUser(t *testing.T) {
	ctx := context.Background()
	svc := usersvc.NewInmemService()

	_, err := svc.GetUser(ctx, "1")
	if err != usersvc.ErrNotFound {
		t.Errorf("did not get ErrNotFound")
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	svc := usersvc.NewInmemService()
	svc.AddUser(ctx, usersvc.User{ID: "1", Name: "bdd"})

	if err := svc.DeleteUser(ctx, "1"); err != nil {
		t.Errorf("failed to delete user: %s", err)
	}

	_, err := svc.GetUser(ctx, "1")
	if err != usersvc.ErrNotFound {
		t.Errorf("user was not deleted")
	}
}

func TestDeleteNonExistentUser(t *testing.T) {
	ctx := context.Background()
	svc := usersvc.NewInmemService()

	if err := svc.DeleteUser(ctx, "1"); err != usersvc.ErrNotFound {
		t.Errorf("should not be able to delete a non-existent user")
	}
}
