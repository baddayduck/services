package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/baddayduck/services/usersvc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	ht "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
)

func main() {
	var (
		consulAddr = flag.String("consul.addr", "", "consul address")
		consulPort = flag.String("consul.port", "", "consul port")
	)
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Service discovery domain.
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()

		consulConfig.Address = fmt.Sprintf("http://%s:%s", *consulAddr, *consulPort)
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	tags := []string{"usersvc"}
	passingOnly := true
	duration := 500 * time.Millisecond
	var userEndpoint endpoint.Endpoint

	ctx := context.Background()
	r := mux.NewRouter()

	factory := userFactory(ctx, "GET", "/users")
	serviceName := "usersvc"
	instancer := consulsd.NewInstancer(client, logger, serviceName, tags, passingOnly)
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(1, duration, balancer)
	userEndpoint = retry

	// GET /sd-users
	r.Methods("GET").Path("/users/{id}").Handler(ht.NewServer(
		userEndpoint,
		decodeConsulGetUserRequest,
		usersvc.EncodeResponse,
	))

	// Interrupt handler.
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP transport.
	go func() {
		logger.Log("transport", "HTTP", "addr", "8080")
		errc <- http.ListenAndServe(":8080", r)
	}()

	// Run!
	logger.Log("exit", <-errc)
}

func userFactory(_ context.Context, method, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = fmt.Sprintf("http://%s", instance)
		}

		tgt, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = path

		var (
			enc ht.EncodeRequestFunc
			dec ht.DecodeResponseFunc
		)
		enc, dec = encodeGetUserRequest, decodeGetUserResponse

		return ht.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil
	}
}

func decodeConsulGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, usersvc.ErrBadRouting
	}
	return usersvc.GetUserRequest{ID: id}, nil
}

func encodeGetUserRequest(_ context.Context, req *http.Request, request interface{}) error {
	lr := request.(usersvc.GetUserRequest)
	p := "/" + lr.ID
	req.URL.Path += p
	return nil
}

func decodeGetUserResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response usersvc.GetUserResponse
	var s map[string]interface{}

	if respCode := resp.StatusCode; respCode >= 400 {
		if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
			return nil, err
		}
		return nil, errors.New(s["error"].(string) + "\n")
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}
