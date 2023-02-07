package errors

import (
	"fmt"
	"runtime"
)

var MaxStackDepth = 50

// Error is an error with stack trace.
type Error interface {
	Error() string
	StackTrace() []Frame
	Unwrap() error
	WithDetails(details interface{}) Error
	Details() interface{}
	Msg() string
}

type errorData struct {
	// err contains original error.
	err error
	// frames contains stack trace of an error.
	frames  []Frame
	msg     string
	details interface{}
}

// WithDetails sets details to error.
func (e *errorData) WithDetails(details interface{}) Error {
	if err, ok := details.(error); ok {
		e.details = err.Error()
	} else {
		e.details = details
	}

	return e
}

// Details return details of error.
func (e *errorData) Details() interface{} {
	return e.details
}

func (e *errorData) WithMsg(msg string) {
	e.msg = msg
}

func (e *errorData) Msg() string {
	if e.msg != "" {
		return e.msg
	}
	return e.err.Error()
}

// FromError creates an error with provided frames.
func FromError(err error, frames []Frame) Error {
	return &errorData{
		err:    err,
		frames: frames,
		msg:    err.Error(),
	}
}

// Errorf creates new error with stacktrace and formatted message.
// Formatting works the same way as in fmt.Errorf.
func Errorf(message string, args ...interface{}) Error {
	return trace(fmt.Errorf(message, args...), defaultCallerSkip)
}

// New creates new error with stacktrace.
func New(message string) Error {
	return trace(fmt.Errorf("%s", message), defaultCallerSkip)
}

// Wrap adds stacktrace to existing error.
func Wrap(err error) Error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if ok {
		return e
	}
	return trace(err, defaultCallerSkip)
}

// Unwrap returns the original error.
func Unwrap(err error) error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if !ok {
		return err
	}
	return e.Unwrap()
}

// Error returns error message.
func (e *errorData) Error() string {
	return e.err.Error()
}

// StackTrace returns stack trace of an error.
func (e *errorData) StackTrace() []Frame {
	return e.frames
}

// Unwrap returns the original error.
func (e *errorData) Unwrap() error {
	return e.err
}

// Frame is a single step in stack trace.
type Frame struct {
	// Func contains a function name.
	Func string
	// Line contains a line number.
	Line int
	// Path contains a file path.
	Path string
}

// StackTrace returns stack trace of an error.
// It will be empty if err is not of type Error.
func StackTrace(err error) []Frame {
	e, ok := err.(Error)
	if !ok {
		return nil
	}
	return e.StackTrace()
}

// String formats Frame to string.
func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s()", f.Path, f.Line, f.Func)
}

func trace(err error, skip int) Error {
	frames := make([]Frame, 0, MaxStackDepth)
	for {
		pc, path, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		frame := Frame{
			Func: fn.Name(),
			Line: line,
			Path: path,
		}
		frames = append(frames, frame)
		skip++
	}
	return &errorData{
		err:    err,
		frames: frames,
	}
}
