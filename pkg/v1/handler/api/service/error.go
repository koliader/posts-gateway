package service

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getErrorCode(err error) codes.Code {
	grpcStatus, _ := status.FromError(err)
	grpcCode := grpcStatus.Code()
	return grpcCode
}
