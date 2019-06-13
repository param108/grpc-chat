package errors

import ()

func NewNotFoundError(err error) GrpcChatError {
	return GrpcChatErrorImpl{e: err, code: NotFoundError}
}
