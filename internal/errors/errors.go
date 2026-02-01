//nolint:revive // Package name matches domain purpose (application errors), stdlib "errors" imported with same name
package errors

import (
	"errors"
	"fmt"
)

// Common sentinel errors.
var (
	// ErrConfigNotFound indicates the configuration file was not found.
	ErrConfigNotFound       = errors.New("fail to find supervisord.conf")
	ErrInvalidStringExpr    = errors.New("invalid string expression format")
	ErrEnvVarNotFound       = errors.New("environment variable not found")
	ErrEnvVarConversion     = errors.New("can't convert environment variable to integer")
	ErrTypeNotImplemented   = errors.New("type not implemented")

	// ErrNoFile indicates a file-related error.
	ErrNoFile               = errors.New("NO_FILE")
	ErrNoExitCode           = errors.New("no exit code")
	ErrEmptyCommand         = errors.New("empty command")
	ErrNoCommandFromString  = errors.New("no command from empty string")
	ErrFailedToGetPID       = errors.New("failed to get pid from file")
	ErrFailedToSetUser      = errors.New("fail to set user")
	ErrProcessNotStarted    = errors.New("process is not started")

	// ErrFailedToReadResult indicates event result reading failure.
	ErrFailedToReadResult   = errors.New("failed to read the result")
	ErrNegativeResultBytes  = errors.New("result bytes is less than 0")

	// ErrNoLog indicates no log is available.
	ErrNoLog                = errors.New("no log")
	ErrNotConnectedToSyslog = errors.New("not connected to syslog server")
	ErrInvalidFormat        = errors.New("invalid format")
	ErrOffsetNegative       = errors.New("offset should not be less than 0")
	ErrLengthNegative       = errors.New("length should not be less than 0")

	// ErrNoSupervisordSection indicates missing supervisord configuration section.
	ErrNoSupervisordSection = errors.New("no supervisord section")
	ErrNegativeValue        = errors.New("negative value")
	ErrNoSuchKey            = errors.New("no such key")
	ErrFailedToGetLimit     = errors.New("fail to get limit")
	ErrLimitExceedsHard     = errors.New("limit exceeds hard limit")
	ErrFailedToSetLimit     = errors.New("fail to set limit")
	ErrBadName              = errors.New("BAD_NAME")
	ErrProcessNotFound      = errors.New("fail to find process")
	ErrInvalidSignalType    = errors.New("signal is not a syscall.Signal")
	ErrNoProcess            = errors.New("no process")
	ErrNotRunning           = errors.New("NOT_RUNNING")

	// ErrBadResponse indicates an invalid XML-RPC response.
	ErrBadResponse          = errors.New("bad response")
	ErrHTTPRequestFailed    = errors.New("fail to send http request to supervisord")
	ErrUnixSocketFailed     = errors.New("fail to connect unix socket path")
	ErrHTTPCreateFailed     = errors.New("fail to create http request")
	ErrUnixSocketWrite      = errors.New("fail to write to unix socket")
	ErrResponseReadFailed   = errors.New("fail to read response")
	ErrIncorrectState       = errors.New("incorrect required state")
)

// NewEnvVarNotFoundError creates an error for missing environment variable.
func NewEnvVarNotFoundError(varName string) error {
	return fmt.Errorf("%w: %s", ErrEnvVarNotFound, varName)
}

// NewEnvVarConversionError creates an error for failed env var conversion.
func NewEnvVarConversionError(varValue string) error {
	return fmt.Errorf("%w: %s", ErrEnvVarConversion, varValue)
}

// NewTypeNotImplementedError creates an error for unimplemented type.
func NewTypeNotImplementedError(typeName string) error {
	return fmt.Errorf("%w: %v", ErrTypeNotImplemented, typeName)
}

// NewNoSuchKeyError creates an error for missing configuration key.
func NewNoSuchKeyError(keyName string) error {
	return fmt.Errorf("%w: %s", ErrNoSuchKey, keyName)
}

// NewNegativeValueError creates an error for negative value.
func NewNegativeValueError(keyName string) error {
	return fmt.Errorf("%w for %s", ErrNegativeValue, keyName)
}

// ErrInvalidArguments is the base error for invalid arguments.
var ErrInvalidArguments = errors.New("invalid arguments")

// NewInvalidArgumentsError creates an error for invalid CLI arguments.
func NewInvalidArgumentsError(usage string) error {
	return fmt.Errorf("%w\nUsage: %s", ErrInvalidArguments, usage)
}

// NewFailedToGetLimitError creates an error for failed limit retrieval.
func NewFailedToGetLimitError(resourceName string) error {
	return fmt.Errorf("%w: %s", ErrFailedToGetLimit, resourceName)
}

// NewLimitExceedsHardError creates an error when limit exceeds hard limit.
func NewLimitExceedsHardError(resourceName string, requested, hardLimit int64) error {
	return fmt.Errorf("%w: %s %d is greater than Hard limit %d", ErrLimitExceedsHard, resourceName, requested, hardLimit)
}

// NewFailedToSetLimitError creates an error for failed limit setting.
func NewFailedToSetLimitError(resourceName string, value int64) error {
	return fmt.Errorf("%w: %s to %d", ErrFailedToSetLimit, resourceName, value)
}

// NewBadNameError creates an error for invalid process name.
func NewBadNameError(processName string) error {
	return fmt.Errorf("%w: no process named %s", ErrBadName, processName)
}

// NewProcessNotFoundError creates an error for process not found.
func NewProcessNotFoundError(processName string) error {
	return fmt.Errorf("%w: %s", ErrProcessNotFound, processName)
}

// NewInvalidSignalTypeError creates an error for invalid signal type.
func NewInvalidSignalTypeError(sigType any) error {
	return fmt.Errorf("%w: %T", ErrInvalidSignalType, sigType)
}

// NewNoProcessError creates an error for no such process.
func NewNoProcessError(processName string) error {
	return fmt.Errorf("%w named %s", ErrNoProcess, processName)
}

// NewBadResponseError creates an error for bad HTTP response.
func NewBadResponseError(statusCode int) error {
	return fmt.Errorf("%w with status code %d", ErrBadResponse, statusCode)
}

// NewHTTPRequestFailedError creates an error for failed HTTP request.
func NewHTTPRequestFailedError(err error) error {
	return fmt.Errorf("%w: %w", ErrHTTPRequestFailed, err)
}

// NewUnixSocketFailedError creates an error for failed Unix socket connection.
func NewUnixSocketFailedError(path string, err error) error {
	return fmt.Errorf("%w: %s: %w", ErrUnixSocketFailed, path, err)
}

// NewHTTPCreateFailedError creates an error for failed HTTP request creation.
func NewHTTPCreateFailedError(err error) error {
	return fmt.Errorf("%w: %w", ErrHTTPCreateFailed, err)
}

// NewUnixSocketWriteError creates an error for failed Unix socket write.
func NewUnixSocketWriteError(path string) error {
	return fmt.Errorf("%w %s", ErrUnixSocketWrite, path)
}

// NewResponseReadFailedError creates an error for failed response read.
func NewResponseReadFailedError(err error) error {
	return fmt.Errorf("%w %w", ErrResponseReadFailed, err)
}
