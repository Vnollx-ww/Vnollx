package utils

import (
	"Vnollx/kitex_gen/base"
)

func WrapWithSuccess(msg string) *base.BaseResponse {
	return &base.BaseResponse{
		Code: 200,
		Msg:  msg,
	}
}
func WrapWithSystemError(msg string) *base.BaseResponse {
	return &base.BaseResponse{
		Code: 500,
		Msg:  msg,
	}
}
func WrapWithLogicError(msg string) *base.BaseResponse {
	return &base.BaseResponse{
		Code: 400,
		Msg:  msg,
	}
}
