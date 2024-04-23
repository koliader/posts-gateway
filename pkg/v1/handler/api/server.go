package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/koliader/posts-gateway/internal/token"
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
	authClient  service.AuthClient
	postsClient service.PostsClient
	tokenMaker  token.Maker
}

func NewServer(config util.Config) (*Server, error) {
	authClient := service.NewAuthClient(config)
	postsClient := service.NewPostsClient(config)
	tokenMaker, err := token.NewJWTMaker(config.TokenKey)
	if err != nil {
		return nil, fmt.Errorf("error to create jwt make: %v", err)
	}
	server := &Server{config: config, authClient: *authClient, postsClient: *postsClient, tokenMaker: tokenMaker}
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
	authRoutes := router.Group("/").Use(s.authMiddleware())

	// auth
	router.POST("/auth/login", s.login)
	router.POST("/auth/register", s.register)

	// users
	router.GET("/users", s.listUsers)
	router.GET("/users/:email", s.getUserByEmail)
	authRoutes.PUT("/users", s.updateUserEmail)

	// posts
	router.GET("/posts/:title", s.getPost)
	router.GET("/posts", s.listPosts)
	router.GET("/posts/byUser/:email", s.listPostsByUser)
	authRoutes.POST("/posts", s.createPost)

	s.router = router
}
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

type Success struct {
	Success bool `json:"success"`
}
