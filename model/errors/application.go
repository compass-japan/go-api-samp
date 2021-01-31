package errors

import "net/http"

/*
 * Application全体とSystem内部で使用するエラーの変数定義
 */

func HttpMethodNotAllowedError(cause error) ApplicationError {
	return &applicationError{
		statusCode:   http.StatusMethodNotAllowed,
		logIgnorable: true,
		systemError: &systemError{
			message: "http method not allowed.",
			cause:   cause,
		},
	}
}

func HttpRouteNotFoundError(cause error) ApplicationError {
	return &applicationError{
		statusCode:   http.StatusNotFound,
		logIgnorable: true,
		systemError: &systemError{
			message: "route not found.",
			cause:   cause,
		},
	}
}

func UnauthorizedError(cause error) ApplicationError {
	return &applicationError{
		statusCode:   http.StatusUnauthorized,
		logIgnorable: true,
		systemError: &systemError{
			message: "unauthorized error.",
			cause:   cause,
		},
	}
}

func InvalidRequestError(cause error) ApplicationError {
	return &applicationError{
		statusCode:   http.StatusBadRequest,
		logIgnorable: false,
		systemError: &systemError{
			message: "invalid request.",
			cause:   cause,
		},
	}
}

func InternalServerError(cause error) ApplicationError {
	return &applicationError{
		statusCode:   http.StatusInternalServerError,
		logIgnorable: false,
		systemError: &systemError{
			message: "internal server error.",
			cause:   cause,
		},
	}
}
