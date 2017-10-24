package authsvc

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware
type Middleware func(Service) Service

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (mw loggingMiddleware) GetSession(ctx context.Context, id string) (sess Session, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetSession", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetSession(ctx, id)
}

func (mw loggingMiddleware) CreateSession(ctx context.Context, sess Session) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "CreateSession", "id", sess.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.CreateSession(ctx, sess)
}

func (mw loggingMiddleware) DeleteSession(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteSession", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteSession(ctx, id)
}
