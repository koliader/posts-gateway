package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koliader/posts-gateway/pkg/v1/handler/api/service"
)

func (s *Server) login(ctx *gin.Context) {
	var req service.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorInvalidArguments(err))
		return
	}

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, code, err := s.auth_client.Login(&c, req)
	if err != nil {
		ctx.JSON(errorCode(*code), errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (s *Server) register(ctx *gin.Context) {
	var req service.RegisterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorInvalidArguments(err))
		return
	}

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, code, err := s.auth_client.Register(&c, req)
	if err != nil {
		ctx.JSON(errorCode(*code), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (s *Server) listUsers(ctx *gin.Context) {

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, code, err := s.auth_client.ListUsers(&c)
	if err != nil {
		ctx.JSON(errorCode(*code), errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res.Users)
}
