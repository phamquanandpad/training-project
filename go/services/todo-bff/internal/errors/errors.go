package errors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"

	"github.com/go-sql-driver/mysql"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
)

const (
	domain   = "owner"
	localeJa = "ja-JP"
)

type ErrorType string

var ErrorTypes = struct {
	AlreadyExistedError     ErrorType
	AuthNError              ErrorType
	AuthZError              ErrorType
	InternalError           ErrorType
	NotFoundError           ErrorType
	ParameterError          ErrorType
	PreconditionFailedError ErrorType
	UnknownError            ErrorType
	CanceledError           ErrorType
}{
	AlreadyExistedError:     "ALREADY_EXISTED_ERROR",
	AuthNError:              "AUTH_N_ERROR",
	AuthZError:              "AUTH_Z_ERROR",
	InternalError:           "INTERNAL_ERROR",
	NotFoundError:           "NOT_FOUND_ERROR",
	ParameterError:          "PARAMETER_ERROR",
	PreconditionFailedError: "PRECONDITIONAL_FAILED_ERROR",
	UnknownError:            "UNKNOWN_ERROR",
	CanceledError:           "CANCELED_ERROR",
}

type Metadata struct {
	key   string
	value string
}

func ToMetadata(key, value string) Metadata {
	return Metadata{key: key, value: value}
}

func ToMetadataInt(key string, value int) Metadata {
	return Metadata{key: key, value: fmt.Sprintf("%d", value)}
}

func ToMetadataInt32(key string, value int32) Metadata {
	return Metadata{key: key, value: fmt.Sprintf("%d", value)}
}

func ToMetadataSlice[T any](key string, values []T) Metadata {
	jsonBytes, err := json.Marshal(&values)
	if err != nil {
		return Metadata{key: key, value: fmt.Sprintf("%v", values)}
	}
	return Metadata{key: key, value: string(jsonBytes)}
}

func WithExecutedPathMetadata() Metadata {
	// notice that we're using 1, so it will actually log where
	// the caller execute this method, 0 = this function, we don't want that.
	_, filename, line, _ := runtime.Caller(1)
	return Metadata{key: "Path", value: fmt.Sprintf("%s:%d", filename, line)}
}

type LocalizedMessage struct {
	JaMessage string
}

type ErrorElement struct {
	Type     ErrorType
	Msg      string
	Err      error
	Metadata map[string]string
}

func NewErrorElement(errtype ErrorType, msg string, err error, mds ...Metadata) ErrorElement {
	md := map[string]string{}
	for _, m := range mds {
		md[m.key] = m.value
	}
	return ErrorElement{
		Type:     errtype,
		Msg:      msg,
		Err:      err,
		Metadata: md,
	}
}

func (e ErrorElement) toError() string {
	md := ""
	if len(e.Metadata) > 0 {
		md = fmt.Sprintf("(metadata: %v)", e.Metadata)
	}
	return fmt.Sprintf("%s: %s %s %v", e.Type, e.Msg, md, e.Err)
}

type AppError struct {
	Elem    ErrorElement
	JaError string
}

func NewAppError(
	errtype ErrorType,
	msg string,
	err error,
	localizedMessage *LocalizedMessage,
	mds ...Metadata,
) AppError {
	jaErrorMessage := ""
	if localizedMessage != nil {
		jaErrorMessage = localizedMessage.JaMessage
	}
	elem := NewErrorElement(errtype, msg, err, mds...)
	return AppError{Elem: elem, JaError: jaErrorMessage}
}

func (e AppError) Error() string {
	return e.Elem.toError()
}

func (e AppError) Unwrap() error {
	return e.Elem.Err
}

func NewAlreadyExistsError(
	msg string,
	err error,
	localizedMessage *LocalizedMessage,
	mds ...Metadata,
) AppError {
	LMessage := &LocalizedMessage{JaMessage: AlreadyExistsJaMessage}
	if localizedMessage != nil {
		LMessage = localizedMessage
	}
	return NewAppError(ErrorTypes.AlreadyExistedError, msg, err, LMessage, mds...)
}

func NewAuthNError(
	msg string,
	err error,
	localizedMessage *LocalizedMessage,
	mds ...Metadata,
) AppError {
	LMessage := &LocalizedMessage{JaMessage: AuthNJaMessage}
	if localizedMessage != nil {
		LMessage = localizedMessage
	}
	return NewAppError(ErrorTypes.AuthNError, msg, err, LMessage, mds...)
}

func NewAuthZError(
	msg string,
	err error,
	localizedMessage *LocalizedMessage,
	mds ...Metadata,
) AppError {
	LMessage := &LocalizedMessage{JaMessage: AuthZJaMessage}
	if localizedMessage != nil {
		LMessage = localizedMessage
	}
	return NewAppError(ErrorTypes.AuthZError, msg, err, LMessage, mds...)
}

func NewNotFoundError(
	msg string,
	err error,
	localizedMessage *LocalizedMessage,
	mds ...Metadata,
) AppError {
	LMessage := &LocalizedMessage{JaMessage: NotFoundJaMessage}
	if localizedMessage != nil {
		LMessage = localizedMessage
	}
	return NewAppError(ErrorTypes.NotFoundError, msg, err, LMessage, mds...)
}

func NewParameterError(
	msg string,
	err error,
	localizedMessage *LocalizedMessage,
	mds ...Metadata,
) AppError {
	LMessage := &LocalizedMessage{JaMessage: PreconditionFailedJaMessage}
	if localizedMessage != nil {
		LMessage = localizedMessage
	}
	return NewAppError(ErrorTypes.ParameterError, msg, err, LMessage, mds...)
}

func NewPreconditionFailedError(
	msg string,
	err error,
	localizedMessage *LocalizedMessage,
	mds ...Metadata,
) AppError {
	LMessage := &LocalizedMessage{JaMessage: InvalidRequestJaMessage}
	if localizedMessage != nil {
		LMessage = localizedMessage
	}
	return NewAppError(ErrorTypes.PreconditionFailedError, msg, err, LMessage, mds...)
}

func NewInternalError(
	msg string,
	err error,
	mds ...Metadata,
) AppError {
	errType := ErrorTypes.InternalError
	// if err is context.Canceled, we should return CanceledError
	if errors.Is(err, context.Canceled) {
		errType = ErrorTypes.CanceledError
	}

	return NewAppError(errType, msg, err, nil, mds...)
}

func NewUnknownError(
	msg string,
	err error,
	mds ...Metadata,
) AppError {
	return NewAppError(ErrorTypes.UnknownError, msg, err, nil, mds...)
}

func NewCanceledError(
	msg string,
	err error,
	mds ...Metadata,
) AppError {
	return NewAppError(ErrorTypes.CanceledError, msg, err, nil, mds...)
}

func ToGRPCCode(err error) codes.Code {
	var appError AppError
	if errors.As(err, &appError) {
		switch appError.Elem.Type {
		case ErrorTypes.AlreadyExistedError:
			return codes.AlreadyExists
		case ErrorTypes.AuthZError:
			return codes.PermissionDenied
		case ErrorTypes.AuthNError:
			return codes.Unauthenticated
		case ErrorTypes.ParameterError:
			return codes.InvalidArgument
		case ErrorTypes.NotFoundError:
			return codes.NotFound
		case ErrorTypes.PreconditionFailedError:
			return codes.FailedPrecondition
		case ErrorTypes.InternalError:
			return codes.Internal
		case ErrorTypes.CanceledError:
			return codes.Canceled
		}
	}

	return codes.Unknown
}

func IsCanceledError(err error) bool {
	return ToGRPCCode(err) == codes.Canceled
}

func IsNotFoundErr(err error) bool {
	if err == nil {
		return false
	}

	var aErr AppError
	if !errors.As(err, &aErr) {
		return false
	}

	return aErr.Elem.Type == ErrorTypes.NotFoundError
}

// GRPCStatus implement the interface { GRPCStatus() *Status } of package grpc/status,
// so that gRPC server can use struct AppError directly when response error
// ref: https://github.com/grpc/grpc-go/blob/v1.65.0/status/status.go#L88-L91
func (e AppError) GRPCStatus() *status.Status {
	stt := status.New(ToGRPCCode(e), e.Elem.Msg)
	errDetails := []protoiface.MessageV1{
		&errdetails.ErrorInfo{
			Reason:   e.Error(),
			Domain:   domain,
			Metadata: e.Elem.Metadata,
		},
	}
	if e.JaError != "" {
		errDetails = append(errDetails, &errdetails.LocalizedMessage{
			Locale:  localeJa,
			Message: e.JaError,
		})
	}

	stt, err := stt.WithDetails(errDetails...)
	if err != nil {
		return status.New(codes.Unknown, fmt.Sprintf("call status.WithDetails failed: %v", err))
	}

	return stt
}

// nolint:exhaustive
func GRPCErrToAppError(err error) AppError {
	grpcError, ok := status.FromError(err)

	if !ok {
		return NewUnknownError("unknown error", err)
	}
	switch grpcError.Code() {
	case codes.Unauthenticated:
		return NewAuthNError(grpcError.Message(), err, nil)
	case codes.PermissionDenied:
		return NewAuthZError(grpcError.Message(), err, nil)
	case codes.FailedPrecondition:
		return NewPreconditionFailedError(grpcError.Message(), err, nil)
	case codes.InvalidArgument:
		return NewParameterError(grpcError.Message(), err, nil)
	case codes.NotFound:
		return NewNotFoundError(grpcError.Message(), err, nil)
	case codes.AlreadyExists:
		return NewAlreadyExistsError(grpcError.Message(), err, nil)
	case codes.Internal:
		return NewInternalError(grpcError.Message(), err)
	default:
		return NewUnknownError(grpcError.Message(), err)
	}
}

func IsMySQLDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062 // MySQL duplicate entry error code
	}
	return false
}
