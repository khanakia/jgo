package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/khanakia/jgo/pkg/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Cli cli.Cli
}

type Server struct {
	Config
	Router       *gin.Engine
	RouterGroups map[string]*gin.RouterGroup
}

func (server Server) AddRouterGroup(name string, path string) *gin.RouterGroup {
	server.RouterGroups[name] = server.Router.Group(path)
	return server.RouterGroups[name]
}

func (server Server) GetRouterGroup(name string) (*gin.RouterGroup, error) {
	group := server.RouterGroups[name]
	if group == nil {
		return nil, errors.New("no route group found")
	}
	return server.RouterGroups[name], nil
}

func (server Server) Start() {
	port := viper.GetString("server.port")
	url := "http://localhost:" + port
	log.Println("Http Sever started at " + color.CyanString(url))
	// log.Println("connect to http://localhost:8082/ for graphql playground")
	server.Router.Run(":" + port)
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func TestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "DDDD",
		})
		c.Abort()
	}
}

// Defining the Graphql handler
// func graphqlHandler() gin.HandlerFunc {
// 	// NewExecutableSchema and Config are in the generated.go file
// 	// Resolver is in the resolver.go file
// 	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

// 	return func(c *gin.Context) {
// 		h.ServeHTTP(c.Writer, c.Request)
// 	}
// }

// // Defining the Playground handler
// func playgroundHandler() gin.HandlerFunc {
// 	h := playground.Handler("GraphQL", "/query")

// 	return func(c *gin.Context) {
// 		h.ServeHTTP(c.Writer, c.Request)
// 	}
// }

func (server Server) CliInit(rootCmd *cobra.Command) {
	var Main = &cobra.Command{
		Use:   "server",
		Short: "Server Pkg - Use `go run . server --help` to see child commands",
		// Long:  `Use go run . greendropship --help`,
		Run: func(cmd *cobra.Command, args []string) {
			color.Yellow("Run below command to see all the child commands.")
			// d := color.New(color.FgMagenta, color.Bold, color.BgHiWhite)
			color.Cyan("go run . server --help")
		},
	}

	var StartCmd = &cobra.Command{
		Use:   "start",
		Short: "start the http server",
		Run: func(cmd *cobra.Command, args []string) {
			server.Start()
		},
	}
	rootCmd.AddCommand(Main)
	Main.AddCommand(StartCmd)
}

func New(config Config) Server {
	// hello.SetFname("aman")
	// fmt.Println(hello.GetFname())

	if viper.GetString("mode") == "production" {
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	// routes(router)

	router.GET("/ping", pingHandler)

	// router.POST("/query", graphqlHandler())
	// router.GET("/gql", playgroundHandler())

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))

	server := Server{
		Config:       config,
		Router:       router,
		RouterGroups: make(map[string]*gin.RouterGroup),
	}

	server.CliInit(server.Config.Cli.RootCmd)

	// authorized1 := router.Group("/p")
	// authorized1.GET("/ping1", pingHandler)

	// authorized := router.Group("/p")
	// authorized := server.AddRouterGroup("private", "/p")
	// authorized.Use(TestMiddleware())
	// authorized.GET("/ping", pingHandler)

	// fmt.Println(authorized)

	// router.Run(viper.GetString("server.port"))

	return server
}
