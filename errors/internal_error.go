package errors

func NewInternalError(err error) GrpcChatError {
	return GrpcChatErrorImpl{e: err, code: InternalError}

}
