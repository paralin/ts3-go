package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/paralin/ts3-go/serverquery"
	"github.com/urfave/cli"
)

var ip string
var username string
var password string

func main() {
	app := cli.NewApp()
	app.Name = "example"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "username",
			Destination: &username,
		},
		cli.StringFlag{
			Name:        "password",
			Destination: &password,
		},
		cli.StringFlag{
			Name:        "ip",
			Destination: &ip,
			Value:       "localhost:10011",
		},
	}
	app.Action = func(c *cli.Context) error {
		if username == "" || password == "" {
			return errors.New("--username and --password must be given")
		}

		ctx := context.Background()
		client, err := serverquery.Dial(ip)
		if err != nil {
			return err
		}

		go client.Run(ctx)
		if err := client.UseServer(ctx, 9987); err != nil {
			return err
		}
		if err := client.Login(ctx, username, password); err != nil {
			return err
		}
		clientList, err := client.GetClientList(ctx)
		if err != nil {
			return err
		}
		dat, _ := json.Marshal(clientList)
		fmt.Printf("client list: %#v\n", string(dat))
		for _, clientSummary := range clientList {
			clientInfo, err := client.GetClientInfo(ctx, clientSummary.Id)
			if err != nil {
				return err
			}
			dat, _ = json.Marshal(clientInfo)
			fmt.Printf("client [%d]: %#v\n", clientSummary.Id, string(dat))
		}
		return nil
	}
	app.RunAndExitOnError()
}
