package gql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/khanakia/jgo/graph"
	"github.com/khanakia/jgo/graph/generated"
	"github.com/khanakia/jgo/pkg/auth"
	"github.com/khanakia/jgo/pkg/server"
)

type Gql struct {
	Config
}

type Config struct {
	Server   server.Server
	Auth     auth.Auth
	Resolver *graph.Resolver
}

// Defining the Graphql handler
func graphqlHandler(resolver *graph.Resolver) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	// h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong11",
	})
}

func New(config Config) Gql {
	// config.Server.Router.GET("/auth/ping", pingHandler)
	// p, err := config.Server.GetRouterGroup("private")
	// fmt.Println(err)
	// if err == nil {
	// 	fmt.Println(p)
	// }
	// p.GET("/a", pingHandler)
	config.Server.Router.POST("/query", graphqlHandler(config.Resolver))
	config.Server.Router.GET("/gql", playgroundHandler())
	// config.Server.Router.GET("/auth/ping11", pingHandler)

	gql := Gql{
		Config: config,
	}

	return gql
}
