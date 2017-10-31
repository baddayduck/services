package usersvc

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

func Register(consulAddress string,
	consulPort string,
	advertiseAddress string,
	advertisePort string) (registrar sd.Registrar) {

	// Logging Domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()
		consulConfig.Address = fmt.Sprintf("%s:%s", consulAddress, consulPort)
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	check := api.AgentServiceCheck{
		HTTP:     fmt.Sprintf("http://%s:%s/health", advertiseAddress, advertisePort),
		Interval: "10s",
		Timeout:  "1s",
		Notes:    "Basic health checks",
	}

	port, _ := strconv.Atoi(advertisePort)
	num := rand.Intn(100)
	asr := api.AgentServiceRegistration{
		ID:      fmt.Sprintf("usersvc-%s", strconv.Itoa(num)), // unique service ID,
		Name:    "usersvc",
		Address: advertiseAddress,
		Port:    port,
		Tags:    []string{"usersvc"},
		Check:   &check,
	}
	return consulsd.NewRegistrar(client, &asr, logger)
}
