package main

import (
	"log"

	"github.com/baddayduck/services/authsvc/db"
	"github.com/baddayduck/services/authsvc/handler"
	proto "github.com/baddayduck/services/authsvc/proto/auth"

	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.auth"),
		micro.Flags(
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "DATABASE_URL",
				Usage:  "The database URL e.g root@tcp(127.0.0.1:3306)/auth",
			},
		),

		micro.Action(func(c *cli.Context) {
			if len(c.String("database_url")) > 0 {
				db.Url = c.String("database_url")
			}
		}),
	)

	service.Init()
	db.Init()

	proto.RegisterAuthHandler(service.Server(), new(handler.Auth))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
