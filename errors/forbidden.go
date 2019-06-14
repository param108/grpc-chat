package errors

func NewForbiddenError(err error) GrpcChatError {
	return GrpcChatErrorImpl{e: err, code: Forbidden}
}
