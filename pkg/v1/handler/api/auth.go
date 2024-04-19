package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koliader/posts-gateway/internal/pb"
	"github.com/koliader/posts-gateway/pkg/v1/handler/api/service"
)

type convertedUser struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
}

func (s *Server) convertUser(user *pb.UserEntity) *convertedUser {
	converted := convertedUser{
		Email:    user.Email,
		Username: user.Username,
	}
	return &converted
}

func (s *Server) login(ctx *gin.Context) {
	var req service.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorInvalidArguments(err))
		return
	}

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, code, err := s.authClient.Login(&c, req)
	if err != nil {
		ctx.JSON(errorCode(code), errorResponse(err))
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

	res, code, err := s.authClient.Register(&c, req)
	if err != nil {
		ctx.JSON(errorCode(code), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (s *Server) listUsers(ctx *gin.Context) {
	var converted []convertedUser
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, code, err := s.authClient.ListUsers(&c)
	if err != nil {
		ctx.JSON(errorCode(code), errorResponse(err))
		return
	}
	for _, user := range res.Users {
		converted = append(converted, *s.convertUser(user))
	}
	ctx.JSON(http.StatusOK, converted)
}

func (s *Server) getUserByEmail(ctx *gin.Context) {

	var req service.GetUserByEmailReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorInvalidArguments(err))
		return
	}
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, code, err := s.authClient.GetUserByEmail(&c, req)
	if err != nil {
		ctx.JSON(errorCode(code), errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, s.convertUser(res.User))
}

type updateUserEmailUriReq struct {
	Email string `uri:"email" binding:"required,email"`
}

type updateUserEmailJsonReq struct {
	Email string `json:"email" binding:"required,email"`
}

func (s *Server) updateUserEmail(ctx *gin.Context) {
	var uriReq updateUserEmailUriReq
	var jsonReq updateUserEmailJsonReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorInvalidArguments(err))
	}
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorInvalidArguments(err))
	}
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req := service.UpdateUserEmailReq{
		Email:   uriReq.Email,
		NewEmil: jsonReq.Email,
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(string)
	headers := service.AuthHeaders{
		Token: authPayload,
	}
	_, code, err := s.authClient.UpdateUserEmail(&c, req, headers)
	if err != nil {
		ctx.JSON(errorCode(code), errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, Success{true})
}
