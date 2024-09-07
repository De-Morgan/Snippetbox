package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golangcollege/sessions"
	"github.com/morgan/snippetbox/pkg/models/database"
)

type contextKey string

const (
	migrationsFolderPath = "./pkg/migration"
	contextKeyUser       = contextKey("user")
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Session Secret")
	dbUser := flag.String("DB_USER", "snippet-admin", "DB user")
	dbPassword := flag.String("DB_PASSWORD", "snippet-admin-password", "DB password")
	dbName := flag.String("DB_NAME", "snippetBox", "DB name")
	dbHost := flag.String("DB_HOST", "mysql", "DB host")

	dsn := fmt.Sprintf("%s:%s@/%s?parseTime=true&multiStatements=true", *dbUser, *dbPassword, *dbName)
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open(*dbHost, dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	// session.Secure = true
	// session.SameSite = http.SameSiteStrictMode

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &database.SnippetModel{DB: db},
		templateCache: templateCache,
		session:       session,
		users:         &database.UserModel{DB: db},
	}
	app.applyMigrations(db, migrationsFolderPath)
	tlsConfig := &tls.Config{PreferServerCipherSuites: true, CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}}
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     app.errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Println("Starting server on ", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
