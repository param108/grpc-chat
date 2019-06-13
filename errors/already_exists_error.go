package errors

func NewAlreadyExistsError(err error) GrpcChatError {
	return GrpcChatErrorImpl{e: err, code: AlreadyExistsError}
}
