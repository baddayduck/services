package usersvc

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware setups up the logging middleware
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) AddUser(ctx context.Context, u User) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "AddUser",
			"id", u.ID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return mw.next.AddUser(ctx, u)
}

func (mw loggingMiddleware) GetUser(ctx context.Context, id string) (p User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetUser",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return mw.next.GetUser(ctx, id)
}

func (mw loggingMiddleware) DeleteUser(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "DeleteUser",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return mw.next.DeleteUser(ctx, id)
}

func (mw loggingMiddleware) HealthCheck() (output bool) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "HealthCheck",
			"result", output,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.HealthCheck()
}
