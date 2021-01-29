package errors

import "net/http"

/*
 * Application全体とSystem内部で使用するエラーの変数定義
 */

var (
	Application = struct {
		HttpMethodNotAllowedError ApplicationErrorBuilder
		HttpRouteNotFoundError    ApplicationErrorBuilder
		UnauthorizedError         ApplicationErrorBuilder
		InvalidRequestError  ApplicationErrorBuilder
		InternalServerError       ApplicationErrorBuilder
	}{
		HttpMethodNotAllowedError: func(cause error, i ...interface{}) ApplicationError {
			return &applicationError{
				statusCode:   http.StatusMethodNotAllowed,
				logIgnorable: true,
				systemError: &systemError{
					message: "http method not allowed.",
					cause:   cause,
				},
			}
		},
		HttpRouteNotFoundError: func(cause error, i ...interface{}) ApplicationError {
			return &applicationError{
				statusCode:   http.StatusNotFound,
				logIgnorable: true,
				systemError: &systemError{
					message: "route not found.",
					cause:   cause,
				},
			}
		},
		UnauthorizedError: func(cause error, i ...interface{}) ApplicationError {
			return &applicationError{
				statusCode:   http.StatusUnauthorized,
				logIgnorable: true,
				systemError: &systemError{
					message: "unauthorized error.",
					cause:   cause,
				},
			}
		},
		InvalidRequestError: func(cause error, i ...interface{}) ApplicationError {
			return &applicationError{
				statusCode:   http.StatusBadRequest,
				logIgnorable: false,
				systemError: &systemError{
					message: "invalid request.",
					cause:   cause,
				},
			}
		},
		InternalServerError: func(cause error, i ...interface{}) ApplicationError {
			return &applicationError{
				statusCode:   http.StatusInternalServerError,
				logIgnorable: false,
				systemError: &systemError{
					message: "internal server error.",
					cause:   cause,
				},
			}
		},
	}

	System = struct {
		UnknownSystemError          SystemErrorBuilder
		DataStoreError              SystemErrorBuilder
		DataStoreValueNotFoundError SystemErrorBuilder
	}{
		UnknownSystemError: func(cause error, i ...interface{}) SystemError {
			return &systemError{
				message: "unknown system error.",
				cause:   cause,
			}
		},
		DataStoreError: func(cause error, i ...interface{}) SystemError {
			return &systemError{
				message: "datastore error.",
				cause:   cause,
			}
		},
		DataStoreValueNotFoundError: func(cause error, i ...interface{}) SystemError {
			return &systemError{
				message: "datastore value not found error.",
				cause:   cause,
			}
		},
	}
)
