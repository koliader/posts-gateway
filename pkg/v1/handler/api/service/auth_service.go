package service

import (
	"context"
	"errors"

	"github.com/koliader/posts-gateway/internal/pb"
	"github.com/koliader/posts-gateway/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// * requests
type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=3"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3"`
}

var (
	authGrpcServiceClient pb.AuthClient
)

// auth gRPC service struct

type AuthClient struct {
	pb.UnimplementedAuthServer
	config util.Config
}

func (ac AuthClient) prepareAuthGrpcClient(c *context.Context) error {
	conn, err := grpc.DialContext(*c, ac.config.AuthGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}...)
	if err != nil {
		authGrpcServiceClient = nil
		return errors.New("connection to auth gRPC service failed")
	}

	// If authGrpcServiceClient already created
	if authGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Success case
	authGrpcServiceClient = pb.NewAuthClient(conn)
	return nil
}

// * this functions calls gRPC service function
func (ac *AuthClient) Register(c *context.Context, req RegisterReq) (*pb.AuthRes, error) {
	// connect auth grpc service
	if err := ac.prepareAuthGrpcClient(c); err != nil {
		return nil, err
	}

	arg := pb.RegisterReq{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}
	// res returns token
	res, err := authGrpcServiceClient.Register(*c, &arg)
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}
	// returning res
	return res, nil
}

func (ac *AuthClient) Login(ctx *context.Context, req LoginReq) (*pb.AuthRes, error) {
	// connect auth grpc service
	if err := ac.prepareAuthGrpcClient(ctx); err != nil {
		return nil, err
	}

	arg := pb.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	}
	res, err := authGrpcServiceClient.Login(*ctx, &arg)
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}
	return res, nil
}
