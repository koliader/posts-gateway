package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/koliader/posts-gateway/internal/util"
	"github.com/koliader/posts-gateway/pkg/v1/handler/api/service"
)

// http rest api server struct
var (
	timeout = time.Second
)

type Server struct {
	config      util.Config
	router      *gin.Engine
	auth_client service.AuthClient
}

func NewServer(config util.Config) (*Server, error) {

	authClient := service.NewAuthClient(config)
	server := &Server{config: config, auth_client: *authClient}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()

	// cors
	c := cors.New(cors.Config{
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
	})

	router.Use(c)

	// auth
	router.POST("/auth/login", s.Login)
	router.POST("/auth/register", s.Register)

	// users
	router.GET("/users", s.ListUsers)
	s.router = router
}
func (s *Server) Start(address string) error {
	// Create a context
	ctx := context.Background()

	// try to connect auth service
	if err := s.auth_client.PrepareAuthGrpcClient(&ctx); err != nil {
		return err
	}

	// Start the server
	return s.router.Run(address)
}
