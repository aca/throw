package throw

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"strings"
)

const maxDepth = 32

type ThrowError struct {
	err        error
	stacktrace []string
}

func (m ThrowError) MarshalJSON() ([]byte, error) {
	v := struct {
		Error      string   `json:"error"`
		Stacktrace []string `json:"stack"`
	}{
		Error:      m.err.Error(),
		Stacktrace: m.stacktrace,
	}
	return json.Marshal(v)
}

func (m ThrowError) Error() string {
	return m.err.Error()
}

func (m ThrowError) Unwrap() error {
	return m.err
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	var a []any
	a = append(a, err)
	return Wrap(fmt.Errorf(format+": %w", a))
}

func Errorf(format string, args ...any) error {
	return Wrap(fmt.Errorf(format, args...))
}

func SlogAttr(err error) slog.Attr {
	return slog.Any("throw", Wrap(err))
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	var terr ThrowError

	// do not re-wrap
	if errors.As(err, &terr) {
		terr.err = err
		return terr
	}

	return ThrowError{err: err, stacktrace: getStackTrace()}
}

func getStackTrace() []string {
	stackBuffer := make([]uintptr, maxDepth)
	length := runtime.Callers(3, stackBuffer[:])
	stack := stackBuffer[:length]

	traceList := make([]string, 0, maxDepth)
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		if goroot != "" && strings.Contains(frame.File, goroot) {
			continue
		}

		if strings.Contains(frame.File, packageName) {
			if !strings.Contains(frame.File, packageName+"/example") {
				continue
			}
		}

		if strings.Contains(frame.File, "try/try.go") {
			continue
		}

		// TODO: add lib to skip trace

		traceList = append(traceList, fmt.Sprintf("%s:%s:%d", frame.Function, frame.File, frame.Line))
	}
	return traceList
}

type fake struct{}

var (
	goroot      = runtime.GOROOT()
	packageName = reflect.TypeOf(fake{}).PkgPath()
)
