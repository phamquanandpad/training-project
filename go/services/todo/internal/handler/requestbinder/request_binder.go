package requestbinder

import (
	"context"
	"encoding/json"

	"github.com/go-playground/validator/v10"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo/internal/errors"
)

type RequestBinder struct {
	validate *validator.Validate
}

func NewRequestBinder(validator *validator.Validate) *RequestBinder {
	return &RequestBinder{
		validate: validator,
	}
}

func (b *RequestBinder) Validate(ctx context.Context, obj interface{}) error {
	if err := b.validate.StructCtx(ctx, obj); err != nil {
		return app_errors.NewParameterError(
			"binder.Validate validate.StructCtx",
			err,
			&app_errors.LocalizedMessage{JaMessage: app_errors.ParameterErrorJaMessage},
		)
	}
	return nil
}

func (b *RequestBinder) Bind(ctx context.Context, req interface{}, dest interface{}) error {
	byteData, err := json.Marshal(req)
	if err != nil {
		return app_errors.NewParameterError(
			"binder.Bind",
			err,
			&app_errors.LocalizedMessage{JaMessage: app_errors.ParameterErrorJaMessage},
		)
	}

	err = json.Unmarshal(byteData, dest)
	if err != nil {
		return app_errors.NewParameterError(
			"binder.Bind",
			err,
			&app_errors.LocalizedMessage{JaMessage: app_errors.ParameterErrorJaMessage},
		)
	}

	if err := b.validate.StructCtx(ctx, dest); err != nil {
		return app_errors.NewParameterError(
			"binder.Bind.validate.StructCtx",
			err,
			&app_errors.LocalizedMessage{JaMessage: app_errors.ParameterErrorJaMessage},
		)
	}

	return nil
}
