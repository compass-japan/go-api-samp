package errors

func UnknownSystemError(cause error) SystemError {
	return &systemError{
		message: "unknown system error.",
		cause:   cause,
	}
}

func DataStoreSystemError(cause error) SystemError {
	return &systemError{
		message: "datastore error.",
		cause:   cause,
	}
}

func DataStoreValueNotFoundSystemError(cause error) SystemError {
	return &systemError{
		message: "datastore value not found error.",
		cause:   cause,
	}
}

func ExAPISystemError(cause error) SystemError {
	return &systemError{
		message: "ex api error.",
		cause:   cause,
	}
}
