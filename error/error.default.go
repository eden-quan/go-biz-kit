package errorutil

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
)

type ErrorsCode struct {
	HttpCode int
	BizCode  int
	Reason   string
}

func NewErrorCode(httpCode int, bizCode int, reason string) *ErrorsCode {
	return &ErrorsCode{
		HttpCode: httpCode,
		BizCode:  bizCode,
		Reason:   reason,
	}
}

func (e ErrorsCode) stackTrace(skipFrame int) string {
	pc := make([]uintptr, 32)
	n := runtime.Callers(skipFrame, pc)
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	msg := make([]string, 0, n)
	for {
		frame, more := frames.Next()
		funcName := frame.Function
		line := frame.Line
		file := frame.File
		msg = append(msg, fmt.Sprintf("\t%s:%d\n\t%s", file, line, funcName))
		if !more {
			break
		}
	}

	return strings.Join(msg, "\n")
}

func (e ErrorsCode) toError(skipFrame int, msgFormat string, args ...interface{}) error {
	err := errors.New(int(e.HttpCode), e.Reason, fmt.Sprintf(msgFormat, args...))
	innerErr := errors.New(int(err.Code), err.Reason, err.Message).WithMetadata(map[string]string{
		"BizCode":        strconv.Itoa(int(e.BizCode)),
		"DefaultMessage": e.Reason,
		"__Stack":        e.stackTrace(skipFrame),
		"__MetaKey":      "__",
	})
	return err.WithCause(innerErr)
}

func (e ErrorsCode) ToError(msgFormat string, args ...interface{}) error {
	return e.toError(4, msgFormat, args...)
}

// FromErrorf generate error from err with extra info, if err is nil, mean's everything is fine, return nil
func (e ErrorsCode) FromErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	te := errors.FromError(e.toError(4, format, args...))
	return te.WithCause(te.Unwrap().(*errors.Error).WithCause(err))
}

// FromError generate error from err with extra info, if err is nil, mean's everything is fine, return nil
func (e ErrorsCode) FromError(err error) error {
	if err == nil {
		return nil
	}

	te := errors.FromError(e.toError(4, ""))
	return te.WithCause(te.Unwrap().(*errors.Error).WithCause(err))
}
