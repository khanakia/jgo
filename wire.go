//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/khanakia/jgo/graph"
	"github.com/khanakia/jgo/pkg/app"
	"github.com/khanakia/jgo/pkg/auth"
	"github.com/khanakia/jgo/pkg/cli"
	"github.com/khanakia/jgo/pkg/dbc"
	"github.com/khanakia/jgo/pkg/gql"
	"github.com/khanakia/jgo/pkg/logger"
	"github.com/khanakia/jgo/pkg/server"
)

func Init() Plugins {
	wire.Build(
		// hello.New,
		cli.New,
		logger.New,
		wire.Struct(new(dbc.Config), "*"),
		dbc.New,
		app.New,
		wire.Struct(new(server.Config), "*"),
		server.New,
		wire.Struct(new(auth.Config), "*"),
		auth.New,
		wire.Struct(new(gql.Config), "*"),
		wire.Struct(new(graph.Resolver), "*"),
		gql.New,
		wire.Struct(new(Plugins), "*"),
	)
	return Plugins{}
}
