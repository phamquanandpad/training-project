package errors

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
}{
	AlreadyExistedError:     "ALREADY_EXISTED_ERROR",
	AuthNError:              "AUTH_N_ERROR",
	AuthZError:              "AUTH_Z_ERROR",
	InternalError:           "INTERNAL_ERROR",
	NotFoundError:           "NOT_FOUND_ERROR",
	ParameterError:          "PARAMETER_ERROR",
	PreconditionFailedError: "PRECONDITIONAL_FAILED_ERROR",
	UnknownError:            "UNKNOWN_ERROR",
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

// nolint: errorlint
func (ae *AppError) UnwrapRootError() error {
	if err, ok := ae.Elem.Err.(AppError); ok {
		return err.UnwrapRootError()
	}

	return ae.Elem.Err
}

// nolint: errorlint
func (ae *AppError) UnwrapRootErrorAsAppError() *AppError {
	if err, ok := ae.Elem.Err.(AppError); ok {
		return err.UnwrapRootErrorAsAppError()
	}

	return ae
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
	LMessage := &LocalizedMessage{JaMessage: InvalidJaMessage}
	return NewAppError(ErrorTypes.InternalError, msg, err, LMessage, mds...)
}

func NewUnknownError(
	msg string,
	err error,
	mds ...Metadata,
) AppError {
	return NewAppError(ErrorTypes.UnknownError, msg, err, nil, mds...)
}

// nolint: exhaustive, cyclop
func GrpcStatusToAppError(err error) AppError {
	grpcError, ok := status.FromError(err)
	if !ok {
		return NewUnknownError("unknown error", err)
	}

	var grpcErrorLocalizedMessage *LocalizedMessage
	grpcErrorMetadata := make([]Metadata, 0)
	for _, detail := range grpcError.Details() {
		switch t := detail.(type) {
		case *errdetails.LocalizedMessage:
			{
				if t.Locale == "ja-JP" {
					grpcErrorLocalizedMessage = &LocalizedMessage{JaMessage: t.Message}
				}
			}
		case *errdetails.ErrorInfo:
			{
				for mdKey, mdValue := range t.Metadata {
					grpcErrorMetadata = append(grpcErrorMetadata, Metadata{
						key:   mdKey,
						value: mdValue,
					})
				}
			}
		}
	}

	switch grpcError.Code() {
	case codes.Unauthenticated:
		return NewAuthNError(grpcError.Message(), err, grpcErrorLocalizedMessage, grpcErrorMetadata...)
	case codes.PermissionDenied:
		return NewAuthZError(grpcError.Message(), err, grpcErrorLocalizedMessage, grpcErrorMetadata...)
	case codes.FailedPrecondition:
		return NewPreconditionFailedError(grpcError.Message(), err, grpcErrorLocalizedMessage, grpcErrorMetadata...)
	case codes.InvalidArgument:
		return NewParameterError(grpcError.Message(), err, grpcErrorLocalizedMessage, grpcErrorMetadata...)
	case codes.NotFound:
		return NewNotFoundError(grpcError.Message(), err, grpcErrorLocalizedMessage, grpcErrorMetadata...)
	case codes.AlreadyExists:
		return NewAlreadyExistsError(grpcError.Message(), err, grpcErrorLocalizedMessage, grpcErrorMetadata...)
	case codes.Internal:
		return NewInternalError(grpcError.Message(), err, grpcErrorMetadata...)
	default:
		return NewUnknownError(grpcError.Message(), err, grpcErrorMetadata...)
	}
}
