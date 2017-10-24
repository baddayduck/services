package authsvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateSessionEndpoint endpoint.Endpoint
	GetSessionEndpoint    endpoint.Endpoint
	DeleteSessionEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateSessionEndpoint: MakeCreateSessionEndpoint(s),
		GetSessionEndpoint:    MakeGetSessionEndpoint(s),
		DeleteSessionEndpoint: MakeDeleteSessionEndpoint(s),
	}
}

func MakeCreateSessionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createSessionRequest)
		e := s.CreateSession(ctx, req.Session)
		return createSessionResponse{Err: e}, nil
	}
}

func MakeGetSessionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getSessionRequest)
		sess, e := s.GetSession(ctx, req.SessionID)
		return getSessionResponse{Session: sess, Err: e}, nil

	}
}

func MakeDeleteSessionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteSessionRequest)
		e := s.DeleteSession(ctx, req.SessionID)
		return deleteSessionResponse{Err: e}, nil
	}
}

type createSessionRequest struct {
	Session Session
}

type createSessionResponse struct {
	Err error `json:"err,omitempty"`
}

func (r createSessionResponse) error() error { return r.Err }

type getSessionRequest struct {
	SessionID string
}

type getSessionResponse struct {
	Session Session `json:"session,omitempty"`
	Err     error   `json:"err,omitempty"`
}

func (r getSessionResponse) error() error { return r.Err }

type deleteSessionRequest struct {
	SessionID string
}

type deleteSessionResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteSessionResponse) error() error { return r.Err }
