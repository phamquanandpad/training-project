package todo

import "context"

type userKey struct{}

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, &userKey{}, user)
}

func ExtractUser(ctx context.Context) *User {
	if v := ctx.Value(&userKey{}); v != nil {
		user, ok := v.(*User)
		if ok {
			return user
		}
	}
	return nil
}
