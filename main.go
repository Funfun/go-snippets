package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Value: "0.0.0.0:8080",
				Usage: "host:port of listen http server",
			},
		},
		Name:  "YourSuperApp",
		Usage: "You are about to start YourSuperApp",
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "your server api",
				Action: func(c *cli.Context) error {
					addr := c.String("addr")
					fmt.Printf("Starting the Test Server at %s\n", addr)

					return http.ListenAndServe(addr, nil)
				},
			},
			{
				Name:    "resource",
				Aliases: []string{"get"},
				Usage:   "your rest api calls",
				Action: func(c *cli.Context) error {
					resource := c.Args().First()

					fmt.Println("Warning: command is no-op")
					fmt.Printf("Requesting resource: %s\n", resource)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
