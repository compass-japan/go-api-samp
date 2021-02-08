package errors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SystemError_Is_Success(t *testing.T) {
	sysErr := UnknownSystemError(errors.New("test error"))
	assert.True(t, errors.Is(sysErr, UnknownSystemError(nil)))
	assert.True(t, sysErr.Is(UnknownSystemError(nil)))
}

func Test_Wrap_Is_Success(t *testing.T) {
	sysErr := UnknownSystemError(errors.New("test error"))
	sysErrNil := UnknownSystemError(nil)
	appErr := InternalServerError(sysErr)
	appErrNil := InternalServerError(sysErrNil)
	assert.True(t, errors.Is(appErr, sysErr))
	assert.True(t, errors.Is(appErrNil, sysErr))
	assert.True(t, errors.Is(appErr, sysErrNil))
	assert.True(t, errors.Is(appErrNil, sysErrNil))
	assert.True(t, errors.Is(appErr, appErrNil))
}
