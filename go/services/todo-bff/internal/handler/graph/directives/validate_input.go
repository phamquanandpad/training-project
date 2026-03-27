package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
)

func ValidateInput(v *validator.Validate) func(ctx context.Context, obj interface{}, next graphql.Resolver, rules *string) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, rules *string) (res interface{}, err error) {
		val, err := next(ctx)
		if err != nil {
			return nil, err
		}

		if rules != nil && *rules != "" {
			// validate Scala type
			if err := v.VarCtx(ctx, val, *rules); err != nil {
				return nil, err
			}
		} else {
			// validate struct type
			// combine with goTag directive to generate `validate` tag on input struct
			if err := v.StructCtx(ctx, val); err != nil {
				return nil, err
			}
		}

		return val, nil
	}

}
