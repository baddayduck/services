package usersvc

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Endpoints collects all of the endpoints that compose a User service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	GetUserEndpoint    endpoint.Endpoint
	AddUserEndpoint    endpoint.Endpoint
	DeleteUserEndpoint endpoint.Endpoint
	HealthEndpoint     endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a Usersvc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetUserEndpoint:    MakeGetUserEndpoint(s),
		AddUserEndpoint:    MakeAddUserEndpoint(s),
		DeleteUserEndpoint: MakeDeleteUserEndpoint(s),
		HealthEndpoint:     MakeHealthEndpoint(s),
	}
}

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in a Usersvc client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	// Note that the request encoders need to modify the request URL, changing
	// the path and method. That's fine: we simply need to provide specific
	// encoders for each endpoint.

	return Endpoints{
		AddUserEndpoint:    httptransport.NewClient("Add", tgt, encodeAddUserRequest, decodeAddUserResponse, options...).Endpoint(),
		GetUserEndpoint:    httptransport.NewClient("GET", tgt, encodeGetUserRequest, decodeGetUserResponse, options...).Endpoint(),
		DeleteUserEndpoint: httptransport.NewClient("DELETE", tgt, encodeDeleteUserRequest, decodeDeleteUserResponse, options...).Endpoint(),
	}, nil
}

// AddUser implements Service. Primarily useful in a client.
func (e Endpoints) AddUser(ctx context.Context, u User) error {
	request := addUserRequest{User: u}
	response, err := e.AddUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(addUserResponse)
	return resp.Err
}

// GetUser implements Service. Primarily useful in a client.
func (e Endpoints) GetUser(ctx context.Context, id string) (User, error) {
	request := GetUserRequest{ID: id}
	response, err := e.GetUserEndpoint(ctx, request)
	if err != nil {
		return User{}, err
	}
	resp := response.(GetUserResponse)
	return resp.User, resp.Err
}

// DeleteUser implements Service. Primarily useful in a client.
func (e Endpoints) DeleteUser(ctx context.Context, id string) error {
	request := deleteUserRequest{ID: id}
	response, err := e.DeleteUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteUserResponse)
	return resp.Err
}

// MakeAddUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeAddUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addUserRequest)
		e := s.AddUser(ctx, req.User)
		return addUserResponse{Err: e}, nil
	}
}

// MakeGetUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserRequest)
		p, e := s.GetUser(ctx, req.ID)
		return GetUserResponse{User: p, Err: e}, nil
	}
}

// MakeDeleteUserEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteUserRequest)
		e := s.DeleteUser(ctx, req.ID)
		return deleteUserResponse{Err: e}, nil
	}
}

// We have two options to return errors from the business logic.
//
// We could return the error via the endpoint itself. That makes certain things
// a little bit easier, like providing non-200 HTTP responses to the client. But
// Go kit assumes that endpoint errors are (or may be treated as)
// transport-domain errors. For example, an endpoint error will count against a
// circuit breaker error count.
//
// Therefore, it's often better to return service (business logic) errors in the
// response object. This means we have to do a bit more work in the HTTP
// response encoder to detect e.g. a not-found error and provide a proper HTTP
// status code. That work is done with the errorer interface, in transport.go.
// Response types that may contain business-logic errors implement that
// interface.

type addUserRequest struct {
	User User
}

type addUserResponse struct {
	Err error `json:"err,omitempty"`
}

func (r addUserResponse) error() error { return r.Err }

type GetUserRequest struct {
	ID string
}

type GetUserResponse struct {
	User User  `json:"User,omitempty"`
	Err  error `json:"err,omitempty"`
}

func (r GetUserResponse) error() error { return r.Err }

type deleteUserRequest struct {
	ID string
}

type deleteUserResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteUserResponse) error() error { return r.Err }

type HealthRequest struct{}
type HealthResponse struct {
	Status bool `json:"status"`
}

func MakeHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		status := s.HealthCheck()
		return HealthResponse{Status: status}, nil
	}
}
