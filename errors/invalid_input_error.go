package errors

func NewInvalidInputError(err error) GrpcChatError {
	return GrpcChatErrorImpl{e: err, code: InvalidInputError}
}
