package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/koliader/posts-gateway/internal/pb"
	"github.com/koliader/posts-gateway/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// * requests
type AuthRes struct {
	Token string `json:"token" binding:"required,jwt"`
}

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
	config util.Config // Add a config field
}

func NewAuthClient(config util.Config) *AuthClient {
	return &AuthClient{
		config: config,
	}
}

func (ac *AuthClient) PrepareAuthGrpcClient(c *context.Context) error {
	conn, err := grpc.DialContext(*c, ac.config.AuthGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}...)
	if err != nil {
		fmt.Println(err)
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
func (ac *AuthClient) Register(ctx *context.Context, req RegisterReq) (res *pb.AuthRes, code *codes.Code, err error) {
	// connect auth grpc service
	if err := ac.PrepareAuthGrpcClient(ctx); err != nil {
		return nil, nil, err
	}

	arg := pb.RegisterReq{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}
	// res returns token
	res, err = authGrpcServiceClient.Register(*ctx, &arg)
	if err != nil {
		grpcStatus, _ := status.FromError(err)
		grpcCode := grpcStatus.Code()

		return nil, &grpcCode, errorResponse(err)
	}
	// returning res
	return res, nil, nil
}

func (ac *AuthClient) Login(ctx *context.Context, req LoginReq) (res *pb.AuthRes, code *codes.Code, err error) {
	// connect auth grpc service
	if err := ac.PrepareAuthGrpcClient(ctx); err != nil {
		return nil, nil, err
	}

	arg := pb.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	}
	res, err = authGrpcServiceClient.Login(*ctx, &arg)
	if err != nil {
		grpcCode := getErrorCode(err)
		return nil, &grpcCode, errorResponse(err)
	}
	return res, nil, nil
}

func (ac *AuthClient) ListUsers(ctx *context.Context) (res *pb.ListUsersRes, code *codes.Code, err error) {
	if err := ac.PrepareAuthGrpcClient(ctx); err != nil {
		return nil, nil, err
	}

	arg := pb.Empty{}
	res, err = authGrpcServiceClient.ListUsers(*ctx, &arg)
	if err != nil {
		grpcCode := getErrorCode(err)
		return nil, &grpcCode, errorResponse(err)
	}
	return res, nil, nil
}
