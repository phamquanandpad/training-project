package middleware

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

const (
	GqlParseFailed                = "GRAPHQL_PARSE_FAILED"
	GqlValidationFailed           = "GRAPHQL_VALIDATION_FAILED"
	GqlBadUserInput               = "BAD_USER_INPUT"
	GqlUnauthenticated            = "UNAUTHENTICATED"
	GqlForbidden                  = "FORBIDDEN"
	GqlPersistedQueryNotFound     = "PERSISTED_QUERY_NOT_FOUND"
	GqlPersistedQueryNotSupported = "PERSISTED_QUERY_NOT_SUPPORTED"
	GqlInternalServerError        = "INTERNAL_SERVER_ERROR"
)

func RecoverPanicError() graphql.RecoverFunc {
	return func(ctx context.Context, recoveredErr interface{}) error {
		var err error

		if panickedStr, ok := recoveredErr.(string); ok {
			err = errors.New(panickedStr)
		}

		if panickedErr, ok := recoveredErr.(error); ok {
			err = panickedErr
		}

		var appErr app_errors.AppError
		if ok := errors.As(err, &appErr); !ok {
			panic(err)
		}

		return err
	}
}

func ErrorPresenter() graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		return convertError(ctx, err)
	}
}

func convertError(ctx context.Context, err error) *gqlerror.Error {

	var validateErr validator.ValidationErrors
	var appErr app_errors.AppError
	// If error caused from validator/v10
	if isValidationErr := errors.As(err, &validateErr); isValidationErr {
		return convertErrorAsAppError(
			ctx,
			app_errors.NewParameterError(
				err.Error(),
				err,
				&app_errors.LocalizedMessage{
					JaMessage: app_errors.InvalidRequestJaMessage,
				},
			),
		)
	}

	// If error casued from internal app error
	if isAppErr := errors.As(err, &appErr); isAppErr {
		return convertErrorAsAppError(ctx, appErr)
	}

	// If cannot parse to any errors
	return convertErrorAsDefaultError(ctx, err)
}

// nolint: cyclop
func convertErrorAsAppError(ctx context.Context, err app_errors.AppError) *gqlerror.Error {
	appErr := err.UnwrapRootErrorAsAppError()
	gqlErrorCode := GqlInternalServerError
	switch appErr.Elem.Type {
	case app_errors.ErrorTypes.AlreadyExistedError:
		{
			gqlErrorCode = GqlBadUserInput
		}
	case app_errors.ErrorTypes.AuthNError:
		{
			gqlErrorCode = GqlUnauthenticated
		}
	case app_errors.ErrorTypes.AuthZError:
		{
			gqlErrorCode = GqlUnauthenticated
		}
	case app_errors.ErrorTypes.InternalError:
		{
			gqlErrorCode = GqlInternalServerError
		}
	case app_errors.ErrorTypes.NotFoundError:
		{
			gqlErrorCode = GqlBadUserInput
		}
	case app_errors.ErrorTypes.ParameterError:
		{
			gqlErrorCode = GqlBadUserInput
		}
	case app_errors.ErrorTypes.PreconditionFailedError:
		{
			gqlErrorCode = GqlBadUserInput
		}
	case app_errors.ErrorTypes.UnknownError:
		{
			gqlErrorCode = GqlInternalServerError
		}
	}

	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: appErr.Elem.Msg,
		Extensions: map[string]interface{}{
			"code":         gqlErrorCode,
			"internalCode": appErr.Elem.Type,
			"jaMessage":    appErr.JaError,
			"metadata":     appErr.Elem.Metadata,
		},
	}
}

// nolint: cyclop
func convertErrorAsDefaultError(ctx context.Context, err error) *gqlerror.Error {
	var gqlErr *gqlerror.Error
	ok := errors.As(err, &gqlErr)

	if !ok {
		return gqlerror.WrapPath(graphql.GetPath(ctx), err)
	}

	var code string
	var internalCode app_errors.ErrorType
	var jaMessage string

	switch gqlErr.Extensions["code"] {
	case GqlParseFailed:
		code = GqlParseFailed
		internalCode = app_errors.ErrorTypes.PreconditionFailedError
		jaMessage = app_errors.PreconditionFailedJaMessage
	case GqlValidationFailed:
		code = GqlValidationFailed
		internalCode = app_errors.ErrorTypes.ParameterError
		jaMessage = app_errors.InvalidRequestJaMessage
	case GqlBadUserInput:
		code = GqlBadUserInput
		internalCode = app_errors.ErrorTypes.ParameterError
		jaMessage = app_errors.InvalidRequestJaMessage
	case GqlUnauthenticated:
		code = GqlUnauthenticated
		internalCode = app_errors.ErrorTypes.AuthNError
		jaMessage = app_errors.AuthNJaMessage
	case GqlForbidden:
		code = GqlForbidden
		internalCode = app_errors.ErrorTypes.AuthZError
		jaMessage = app_errors.AuthZJaMessage
	case GqlPersistedQueryNotFound:
		code = GqlPersistedQueryNotFound
		internalCode = app_errors.ErrorTypes.InternalError
		jaMessage = app_errors.InvalidJaMessage
	case GqlPersistedQueryNotSupported:
		code = GqlPersistedQueryNotSupported
		internalCode = app_errors.ErrorTypes.InternalError
		jaMessage = app_errors.InvalidJaMessage
	case GqlInternalServerError:
		code = GqlInternalServerError
		internalCode = app_errors.ErrorTypes.InternalError
		jaMessage = app_errors.InvalidJaMessage
	default:
		code = GqlInternalServerError
		internalCode = app_errors.ErrorTypes.UnknownError
		jaMessage = app_errors.InvalidJaMessage
	}
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: gqlErr.Message,
		Extensions: map[string]interface{}{
			"code":         code,
			"internalCode": internalCode,
			"jaMessage":    jaMessage,
		},
	}
}
