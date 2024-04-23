package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/koliader/posts-gateway/internal/pb"
	"github.com/koliader/posts-gateway/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// * headers
type TokenHeader struct {
	Email string
}

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

type GetUserByEmailReq struct {
	Email string `uri:"email" binding:"required,email"`
}

type UpdateUserEmailReq struct {
	NewEmil string
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
		code := grpcStatus.Code()

		return nil, &code, errorResponse(err)
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
		code := getErrorCode(err)
		return nil, &code, errorResponse(err)
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
		code := getErrorCode(err)
		return nil, &code, errorResponse(err)
	}
	return res, nil, nil
}

func (ac *AuthClient) GetUserByEmail(ctx *context.Context, req GetUserByEmailReq) (res *pb.UserRes, code *codes.Code, err error) {
	if err := ac.PrepareAuthGrpcClient(ctx); err != nil {
		return nil, nil, err
	}

	arg := pb.GetUserByEmailReq{
		Email: req.Email,
	}
	res, err = authGrpcServiceClient.GetUserByEmail(*ctx, &arg)
	if err != nil {
		code := getErrorCode(err)
		return nil, &code, errorResponse(err)
	}
	return res, nil, nil
}

func (ac *AuthClient) UpdateUserEmail(ctx *context.Context, req UpdateUserEmailReq, header TokenHeader) (res *pb.UserRes, codes *codes.Code, err error) {
	md := metadata.New(map[string]string{
		"authorization": header.Email,
	})

	// Attach metadata to context
	ctxWithMetadata := metadata.NewOutgoingContext(*ctx, md)

	if err := ac.PrepareAuthGrpcClient(ctx); err != nil {
		return nil, nil, err
	}
	arg := pb.UpdateUserEmailReq{
		NewEmail: req.NewEmil,
	}
	res, err = authGrpcServiceClient.UpdateUserEmail(ctxWithMetadata, &arg)
	if err != nil {
		code := getErrorCode(err)
		return nil, &code, errorResponse(err)
	}
	return res, nil, nil
}
