package datastore

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/phamquanandpad/training-project/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/services/todo/internal/errors"
)

type TodoDBKey struct{}

func WithTodoDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, &TodoDBKey{}, db)
}

func ExtractTodoDB(ctx context.Context) (*gorm.DB, error) {
	if v := ctx.Value(&TodoDBKey{}); v != nil {
		db, ok := v.(*gorm.DB)
		if ok {
			return db, nil
		}
	}
	return nil, errors.NewInternalError("ExtractTodoDB: failed to extract DB", nil)
}

func WithCursorPagingTokenWhereScope(columnFields []CursorPagingField, token *string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if token == nil || len(*token) == 0 {
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

func WithCursorPagingTimeWhereScope(column string, pagingParam todo.CursorPagingParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagingParam.LastDataAt == nil {
			return db
		}

		// Get newer data
		if pagingParam.SortingOrder == todo.SortingOrders.Asc {
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
		if pagingParam.SortingOrder == todo.SortingOrders.Desc {
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

func WithCursorPagingTimeGroupHavingScope(column string, pagingParam todo.CursorPagingParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagingParam.LastDataAt == nil {
			return db
		}

		// Get newer data
		if pagingParam.SortingOrder == todo.SortingOrders.Asc {
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
		if pagingParam.SortingOrder == todo.SortingOrders.Desc {
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

func WithCursorPagingSortableIDWhereScope(idColumn string, pagingParam todo.CursorPagingParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagingParam.Token == nil {
			return db
		}

		// Get newer data
		if pagingParam.SortingOrder == todo.SortingOrders.Asc {
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
		if pagingParam.SortingOrder == todo.SortingOrders.Desc {
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
