package errors

import "fmt"

/*
 *  エラーの構造定義
 *  getterと標準エラー(errors package)のラップ
 */

type (
	SystemError interface {
		error
		Is(error) bool
		Unwrap() error
		Message() string
		Causes() []string
	}

	ApplicationError interface {
		SystemError
		StatusCode() int
		LogIgnorable() bool
	}

	systemError struct {
		message string
		cause   error
	}

	applicationError struct {
		statusCode   int
		logIgnorable bool
		*systemError
	}
)

func (e *systemError) Error() string {
	return e.message
}

func (e *systemError) Is(target error) bool {
	return e.Error() == target.Error()
}

func (e *systemError) Unwrap() error {
	return e.cause
}

func (e *systemError) Message() string {
	return e.message
}

func (e *systemError) Causes() []string {
	causes := []string{e.Error()}

	cur := e.Unwrap()
	for cur != nil {
		causes = append(causes, cur.Error())
		e, ok := cur.(interface {
			error
			Unwrap() error
		})
		if !ok {
			break
		}
		cur = e.Unwrap()
	}

	return causes
}

func (e *applicationError) Error() string {
	return fmt.Sprintf("status: %d, e: %s", e.statusCode, e.message)
}

func (e *applicationError) StatusCode() int {
	return e.statusCode
}

func (e *applicationError) LogIgnorable() bool {
	return e.logIgnorable
}
