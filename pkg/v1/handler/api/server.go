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
	timeout     = time.Second
	auth_client service.AuthClient
)

type Server struct {
	config util.Config
	router *gin.Engine
}

func NewServer(config util.Config) (*Server, error) {

	server := &Server{config: config}

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
	s.router = router
}
func (s *Server) Start(address string) error {
	// Create a context
	ctx := context.Background()

	// try to connect auth service
	authClient := service.AuthClient{}
	if err := authClient.PrepareAuthGrpcClient(&ctx); err != nil {
		return err
	}

	// Start the server
	return s.router.Run(address)
}
