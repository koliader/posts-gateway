package service

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getErrorCode(err error) codes.Code {
	grpcStatus, _ := status.FromError(err)
	grpcCode := grpcStatus.Code()
	return grpcCode
}

func errorResponse(err error) error {
	return errors.New(status.Convert(err).Message())
}
