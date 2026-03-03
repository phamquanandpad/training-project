package datastore

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/errors"
)

type AuthDBKey struct{}

func WithAuthDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, &AuthDBKey{}, db)
}

func ExtractAuthDB(ctx context.Context) (*gorm.DB, error) {
	if v := ctx.Value(&AuthDBKey{}); v != nil {
		db, ok := v.(*gorm.DB)
		if ok {
			return db, nil
		}
	}
	return nil, errors.NewInternalError("ExtractAuthDB: failed to extract DB", nil)
}

func WithCursorPagingTokenWhereScope(columnFields []CursorPagingField, token *string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(*token) == 0 {
			return db
		}

		values := ParsePageToken(*token)
		if len(values) != len(columnFields) {
			return db
		}

		columnValues := make([]any, len(values))
		for idx, value := range values {
			columnValues[idx] = value
		}

		sql, args := BuildCursorPagingCondition(columnFields, columnValues)
		if sql == "" {
			return db
		}

		return db.Where(sql, args...)
	}
}

func WithCursorPagingTimeWhereScope(column string, pagingParam auth_models.CursorPagingParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagingParam.LastDataAt == nil {
			return db
		}

		// Get newer data
		if pagingParam.SortingOrder == auth_models.SortingOrders.Asc {
			if pagingParam.HasCursor {
				return db.Where(
					fmt.Sprintf("%s >= ?", column),
					pagingParam.LastDataAt,
				)
			}
			return db.Where(
				fmt.Sprintf("%s > ?", column),
				pagingParam.LastDataAt,
			)
		}

		// Get older data
		if pagingParam.SortingOrder == auth_models.SortingOrders.Desc {
			if pagingParam.HasCursor {
				return db.Where(
					fmt.Sprintf("%s <= ?", column),
					pagingParam.LastDataAt,
				)
			}

			return db.Where(
				fmt.Sprintf("%s < ?", column),
				pagingParam.LastDataAt,
			)
		}

		return db
	}
}

func WithCursorPagingTimeGroupHavingScope(column string, pagingParam auth_models.CursorPagingParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagingParam.LastDataAt == nil {
			return db
		}

		// Get newer data
		if pagingParam.SortingOrder == auth_models.SortingOrders.Asc {
			if pagingParam.HasCursor {
				return db.Having(
					fmt.Sprintf("%s >= ?", column),
					pagingParam.LastDataAt,
				)
			}
			return db.Having(
				fmt.Sprintf("%s > ?", column),
				pagingParam.LastDataAt,
			)
		}

		// Get older data
		if pagingParam.SortingOrder == auth_models.SortingOrders.Desc {
			if pagingParam.HasCursor {
				return db.Having(
					fmt.Sprintf("%s <= ?", column),
					pagingParam.LastDataAt,
				)
			}

			return db.Having(
				fmt.Sprintf("%s < ?", column),
				pagingParam.LastDataAt,
			)
		}

		return db
	}
}

func WithCursorPagingSortableIDWhereScope(idColumn string, pagingParam auth_models.CursorPagingParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagingParam.Token == nil {
			return db
		}

		// Get newer data
		if pagingParam.SortingOrder == auth_models.SortingOrders.Asc {
			if pagingParam.HasCursor {
				return db.Where(
					fmt.Sprintf("%s >= ?", idColumn),
					pagingParam.Token,
				)
			}
			return db.Where(
				fmt.Sprintf("%s > ?", idColumn),
				pagingParam.Token,
			)
		}

		// Get older data
		if pagingParam.SortingOrder == auth_models.SortingOrders.Desc {
			if pagingParam.HasCursor {
				return db.Where(
					fmt.Sprintf("%s <= ?", idColumn),
					pagingParam.Token,
				)
			}

			return db.Where(
				fmt.Sprintf("%s < ?", idColumn),
				pagingParam.Token,
			)
		}

		return db
	}
}
