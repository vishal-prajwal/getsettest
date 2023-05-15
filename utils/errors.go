package utils

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

type NewGrpcError struct {
	Code codes.Code
	Err  error
}

func (r *NewGrpcError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.Code, r.Err)
}
