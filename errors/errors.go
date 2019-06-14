package errors

import ()

const (
	NotFoundError      = "not-found"
	InternalError      = "internal-error"
	AlreadyExistsError = "already-exists"
	InvalidInputError  = "invalid-input"
	Forbidden          = "Forbidden"
)

type GrpcChatError interface {
	error
	Code() string
}

type GrpcChatErrorImpl struct {
	e    error
	code string
}

func (g GrpcChatErrorImpl) Error() string {
	return g.e.Error()
}

func (g GrpcChatErrorImpl) Code() string {
	return g.code
}
