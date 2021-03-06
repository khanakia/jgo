// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/khanakia/jgo/graph"
	"github.com/khanakia/jgo/pkg/app"
	"github.com/khanakia/jgo/pkg/auth"
	"github.com/khanakia/jgo/pkg/cli"
	"github.com/khanakia/jgo/pkg/dbc"
	"github.com/khanakia/jgo/pkg/gql"
	"github.com/khanakia/jgo/pkg/logger"
	"github.com/khanakia/jgo/pkg/server"
)

// Injectors from wire.go:

func Init() Plugins {
	cliCli := cli.New()
	appApp := app.New()
	loggerLogger := logger.New()
	config := dbc.Config{
		Logger: loggerLogger,
	}
	dbcDbc := dbc.New(config)
	serverConfig := server.Config{
		Cli: cliCli,
	}
	serverServer := server.New(serverConfig)
	authConfig := auth.Config{
		Dbc:    dbcDbc,
		Logger: loggerLogger,
		Server: serverServer,
	}
	authAuth := auth.New(authConfig)
	resolver := &graph.Resolver{
		Auth: authAuth,
	}
	gqlConfig := gql.Config{
		Server:   serverServer,
		Auth:     authAuth,
		Resolver: resolver,
	}
	gqlGql := gql.New(gqlConfig)
	plugins := Plugins{
		Cli:    cliCli,
		App:    appApp,
		Db:     dbcDbc,
		Auth:   authAuth,
		Server: serverServer,
		Gql:    gqlGql,
	}
	return plugins
}
