package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koliader/posts-gateway/pkg/v1/handler/api/service"
)

func (s *Server) createPost(ctx *gin.Context) {
	var req service.CreatePostReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorInvalidArguments(err))
		return
	}
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, code, err := s.posts_client.CreatePost(&c, req)
	if err != nil {
		ctx.JSON(errorCode(code), errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res.Post)
}

func (s *Server) getPost(ctx *gin.Context) {
	var req service.GetPostReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorInvalidArguments(err))
		return
	}

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, code, err := s.posts_client.GetPost(&c, req)
	if err != nil {
		ctx.JSON(errorCode(code), errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res.Post)
}

func (s *Server) listPosts(ctx *gin.Context) {
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, code, err := s.posts_client.ListPosts(&c)
	if err != nil {
		ctx.JSON(errorCode(code), errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res.Posts)
}

func (s *Server) listPostsByUser(ctx *gin.Context) {
	var req service.ListPostsByUserReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorInvalidArguments(err))
		return
	}
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, code, err := s.posts_client.ListPostsByUser(&c, req)

	if err != nil {
		ctx.JSON(errorCode(code), errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res.Posts)
}
