package context

import (
	"context"
	"time"

	"github.com/phamquanandpad/training-project/services/todo/internal/errors"
)

// ValueOnly returns a context that keeps only the parent's values.
func ValueOnly(ctx context.Context) context.Context {
	return valueOnlyContext{ctx}
}

type valueOnlyContext struct {
	context.Context
}

func (v valueOnlyContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (v valueOnlyContext) Done() <-chan struct{} {
	return nil
}

func (v valueOnlyContext) Err() error {
	return nil
}

type serviceNameKey struct{}

func WithServiceName(ctx context.Context, serviceName string) context.Context {
	return context.WithValue(ctx, &serviceNameKey{}, serviceName)
}

func ExtractServiceName(ctx context.Context) (string, error) {
	if v := ctx.Value(&serviceNameKey{}); v != nil {
		s, ok := v.(string)
		if ok {
			return s, nil
		}
	}
	return "", errors.NewInternalError("ExtractServiceName: failed to extract ServiceName", nil)
}

type methodNameKey struct{}

func WithMethodName(ctx context.Context, methodName string) context.Context {
	return context.WithValue(ctx, &methodNameKey{}, methodName)
}

func ExtractMethodName(ctx context.Context) (string, error) {
	if v := ctx.Value(&methodNameKey{}); v != nil {
		s, ok := v.(string)
		if ok {
			return s, nil
		}
	}
	return "", errors.NewInternalError("ExtractMethodName: failed to method Name", nil)
}
