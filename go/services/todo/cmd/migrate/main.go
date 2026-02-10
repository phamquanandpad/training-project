package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/phamquanandpad/training-project/services/todo/internal/config"

	"github.com/golang-migrate/migrate/v4"
	mysqlDriver "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/go-sql-driver/mysql"
)

func buildDSN(c config.DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		c.DBUser, c.DBPass, c.DBHost, strconv.Itoa(c.DBPort), c.DBName,
	)
}

func main() {
	cfg, err := config.LoadDBConfig()
	if err != nil {
		log.Fatal(err)
	}

	dsn := buildDSN(*cfg)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := mysqlDriver.WithInstance(db, &mysqlDriver.Config{})
	if err != nil {
		log.Fatal(err)
	}

	wd, _ := os.Getwd()
	fmt.Println(wd)

	path := filepath.Join(wd, "database", "migrations")

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+path,
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	action := "up"
	steps := 0

	if len(os.Args) > 1 {
		action = os.Args[1]
	}
	if len(os.Args) > 2 {
		steps, _ = strconv.Atoi(os.Args[2])
	}

	switch action {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "steps":
		err = m.Steps(steps)
	case "force":
		err = m.Force(steps)
	case "version":
		v, dirty, _ := m.Version()
		fmt.Println("version:", v, "dirty:", dirty)
		return
	default:
		log.Fatalf("unknown command: %s", action)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	fmt.Println("Done")
}
