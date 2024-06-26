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
)

var (
	postsGrpcServiceClient pb.PostClient
)

// * headers

type CreatePostReq struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

type GetPostReq struct {
	Title string `uri:"title" binding:"required"`
}

type ListPostsByUserReq struct {
	Email string `uri:"email" binding:"required,email"`
}

type PostsClient struct {
	pb.UnimplementedPostServer
	config util.Config
}

func NewPostsClient(config util.Config) *PostsClient {
	return &PostsClient{
		config: config,
	}
}

func (pc *PostsClient) PreparePostsGrpcClient(c *context.Context) error {
	fmt.Println(pc.config.PostsGrpcService)
	fmt.Println(pc.config.AuthGrpcService)
	conn, err := grpc.DialContext(*c, pc.config.PostsGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}...)
	if err != nil {
		fmt.Println(err)
		postsGrpcServiceClient = nil
		return errors.New("connection to posts gRPC service failed")
	}

	// If postsGrpcServiceClient already created
	if postsGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Success case
	postsGrpcServiceClient = pb.NewPostClient(conn)
	return nil
}

func (pc *PostsClient) CreatePost(c *context.Context, req CreatePostReq, header TokenHeader) (res *pb.CreatePostRes, code *codes.Code, err error) {
	md := metadata.New(map[string]string{
		"authorization": header.Email,
	})

	ctxWithMetadata := metadata.NewOutgoingContext(*c, md)

	if err := pc.PreparePostsGrpcClient(c); err != nil {
		return nil, nil, err
	}
	arg := pb.CreatePostReq{
		Title: req.Title,
		Body:  req.Body,
	}

	res, err = postsGrpcServiceClient.CreatePost(ctxWithMetadata, &arg)
	if err != nil {
		code := getErrorCode(err)
		return nil, &code, errorResponse(err)
	}
	return res, nil, nil
}

func (pc *PostsClient) GetPost(c *context.Context, req GetPostReq) (res *pb.GetPostRes, code *codes.Code, err error) {
	if err := pc.PreparePostsGrpcClient(c); err != nil {
		return nil, nil, err
	}

	arg := pb.GetPostReq{
		Title: req.Title,
	}
	res, err = postsGrpcServiceClient.GetPost(*c, &arg)
	if err != nil {
		code := getErrorCode(err)
		return nil, &code, errorResponse(err)
	}
	return res, nil, nil
}

func (pc *PostsClient) ListPosts(c *context.Context) (res *pb.ListPostsRes, code *codes.Code, err error) {
	if err := pc.PreparePostsGrpcClient(c); err != nil {
		return nil, nil, err
	}
	res, err = postsGrpcServiceClient.ListPosts(*c, &pb.Empty{})
	if err != nil {
		code := getErrorCode(err)
		return nil, &code, errorResponse(err)
	}

	return res, nil, nil
}

func (pc *PostsClient) ListPostsByUser(c *context.Context, req ListPostsByUserReq) (res *pb.ListPostsRes, code *codes.Code, err error) {
	if err := pc.PreparePostsGrpcClient(c); err != nil {
		return nil, nil, err
	}
	arg := pb.ListPostsByUserReq{
		Owner: req.Email,
	}
	res, err = postsGrpcServiceClient.ListPostsByUser(*c, &arg)
	if err != nil {
		code := getErrorCode(err)
		return nil, &code, errorResponse(err)
	}
	return res, nil, nil
}
