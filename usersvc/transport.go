package usersvc

// The Usersvc is just over HTTP, so we just have a single transport.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a Usersvc server.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// GET		/health			basic health check
	// POST    	/users			adds another User
	// GET     	/users/:id		retrieves the given User by id
	// DELETE  	/users/:id		remove the given User

	r.Methods("GET").Path("/health").Handler(httptransport.NewServer(
		e.HealthEndpoint,
		decodeHealthRequest,
		EncodeResponse,
		options...,
	))

	r.Methods("POST").Path("/users").Handler(httptransport.NewServer(
		e.AddUserEndpoint,
		decodeaddUserRequest,
		EncodeResponse,
		options...,
	))

	r.Methods("GET").Path("/users/{id}").Handler(httptransport.NewServer(
		e.GetUserEndpoint,
		decodeGetUserRequest,
		EncodeResponse,
		options...,
	))

	r.Methods("DELETE").Path("/users/{id}").Handler(httptransport.NewServer(
		e.DeleteUserEndpoint,
		decodeDeleteUserRequest,
		EncodeResponse,
		options...,
	))
	return r
}

func decodeaddUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req addUserRequest
	if e := json.NewDecoder(r.Body).Decode(&req.User); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return GetUserRequest{ID: id}, nil
}

func decodeDeleteUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteUserRequest{ID: id}, nil
}

func encodeAddUserRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/users/")
	req.Method, req.URL.Path = "POST", "/users/"
	return encodeRequest(ctx, req, request)
}

func encodeGetUserRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/users/{id}")
	r := request.(GetUserRequest)
	UserID := url.QueryEscape(r.ID)
	req.Method, req.URL.Path = "GET", "/users/"+UserID
	return encodeRequest(ctx, req, request)
}

func encodeDeleteUserRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/users/{id}")
	r := request.(deleteUserRequest)
	UserID := url.QueryEscape(r.ID)
	req.Method, req.URL.Path = "DELETE", "/users/"+UserID
	return encodeRequest(ctx, req, request)
}

func decodeAddUserResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response addUserResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetUserResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response GetUserResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteUserResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteUserResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeHealthRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return HealthRequest{}, nil
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// EncodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// Usersvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
