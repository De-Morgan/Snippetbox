package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golangcollege/sessions"
	"github.com/morgan/snippetbox/pkg/models/database"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *database.SnippetModel
	templateCache map[string]*template.Template
	session       *sessions.Session
	users         *database.UserModel
}

func (app *application) applyMigrations(db *sql.DB, migrationsFolderPath string) error {
	app.infoLog.Println("Applying database migrations...")

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsFolderPath),
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migration apply error: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up error: %w", err)
	}

	return nil
}
