package main

import (
	"log"
	"os"
	"time"

	"github.com/XuVic/miniserver/client"
	"github.com/XuVic/miniserver/handler"
	"github.com/XuVic/miniserver/server"
	"github.com/urfave/cli"
)

var app = cli.NewApp()

var addr string

func main() {

	info()
	commands()

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func info() {
	app.Name = "MiniServer"
	app.Author = "VicXu"
	app.Version = "1.0.0"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "Running a Miniserver.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "addr, a", Value: "localhost"},
				cli.StringFlag{Name: "port, p", Value: "3000"},
			},
			Action: func(c *cli.Context) {
				mux := handler.Routes()
				engine := server.NewEngine(c.String("addr"), c.String("port"))
				mux.HandleFunc("/stat", engine.HandleStat)
				engine.Handler = mux
				engine.TimeOut = time.Second * 10
				engine.RunTCP()
			},
		}, {
			Name:  "client",
			Usage: "Running a Client.",
			Action: func(c *cli.Context) {
				client.RunPrompter()
			},
		},
	}
}
