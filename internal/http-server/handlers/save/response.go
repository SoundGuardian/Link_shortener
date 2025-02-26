package save

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func OK() Response {
	return Response{
		Status: StatusOk,
	}
}
