package main

//go:generate go run ./gqlgen.go

import (
	"github.com/khanakia/jgo/pkg/app"
	"github.com/khanakia/jgo/pkg/auth"
	"github.com/khanakia/jgo/pkg/cli"
	"github.com/khanakia/jgo/pkg/dbc"
	"github.com/khanakia/jgo/pkg/gql"
	"github.com/khanakia/jgo/pkg/server"
)

type Plugins struct {
	Cli    cli.Cli
	App    app.App
	Db     dbc.Dbc
	Auth   auth.Auth
	Server server.Server
	Gql    gql.Gql
}

func main() {
	p := Init()
	// fmt.Println(p.App.Version())
	p.Cli.Execute()
	// p.Server.Start()
}
