package datastore

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/config"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
)

type TodoConn struct {
	GormDB *gorm.DB
}

func NewTodoSQLHandler(conf *config.DBConfig) (*TodoConn, func(), error) {
	sourceDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		conf.DBUser,
		conf.DBPass,
		conf.DBHost,
		conf.DBPort,
		conf.DBName,
	)

	db, closer, err := newSQLHandler(sourceDSN)
	return &TodoConn{GormDB: db}, closer, err
}

func newSQLHandler(dsn string) (*gorm.DB, func(), error) {
	if _, err := time.LoadLocation("UTC"); err != nil {
		return nil, nil, fmt.Errorf(": %w", err)
	}

	// nolint: exhaustivestruct
	conn, err := gorm.Open(
		mysql.New(mysql.Config{
			DriverName: "mysql",
			DSN:        dsn,
		}),
		&gorm.Config{
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf(": %w", err)
	}
	err = conn.Use(dbresolver.Register(dbresolver.Config{
		Sources:           []gorm.Dialector{mysql.Open(dsn)},
		TraceResolverMode: true,
	}))
	if err != nil {
		return nil, nil, fmt.Errorf(": %w", err)
	}
	db, err := conn.DB()
	if err != nil {
		return nil, nil, fmt.Errorf(": %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf(": %w", err)
	}

	conn.Set("gorm:table_options", "ENGINE=InnoDB")

	return conn, func() {
		_ = db.Close()
	}, nil
}

type binder struct {
	todoConn *TodoConn
}

func NewConnectionBinder(todoConn *TodoConn) gateway.Binder {
	return &binder{
		todoConn: todoConn,
	}
}

func (b binder) Bind(ctx context.Context) context.Context {
	return WithTodoDB(ctx, b.todoConn.GormDB)
}
