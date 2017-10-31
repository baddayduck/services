package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/baddayduck/services/authsvc"
	"github.com/go-kit/kit/log"
)

func main() {
	var (
		consulAddr    = flag.String("consul.addr", "", "consul address")
		consulPort    = flag.String("consul.port", "", "consul port")
		advertiseAddr = flag.String("advertise.addr", "", "advertise address")
		advertisePort = flag.String("advertise.port", "", "advertise port")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s authsvc.Service
	{
		s = authsvc.NewInmemService()
		s = authsvc.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = authsvc.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	registrar := authsvc.Register(
		*consulAddr,
		*consulPort,
		*advertiseAddr,
		*advertisePort,
	)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "port", *advertisePort)
		// register service
		registrar.Register()
		errs <- http.ListenAndServe(fmt.Sprintf(":%s", *advertisePort), h)
	}()
	err := <-errs
	// deregister service
	registrar.Deregister()
	logger.Log("exit", err)
}
