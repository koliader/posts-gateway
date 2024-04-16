package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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
	router.POST("/auth/login", s.login)
	router.POST("/auth/register", s.register)

	// users
	router.GET("/users", s.listUsers)
	s.router = router
}
func (s *Server) Start(address string) error {
	ctx := context.Background()

	if err := s.auth_client.PrepareAuthGrpcClient(&ctx); err != nil {
		log.Error().Err(err).Msg("failed to connect to auth service")
		return err
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("auth service connected")

	return s.router.Run(address)
}
