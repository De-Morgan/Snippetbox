package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/morgan/snippetbox/pkg/models/database"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "root:secret@/snippets?parseTime=true&multiStatements=true", "MYSQL DSN")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open("mysql", *dsn)
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

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &database.SnippetModel{DB: db},
		templateCache: templateCache,
	}
	srv := &http.Server{Addr: *addr, ErrorLog: app.errorLog, Handler: app.routes()}
	infoLog.Println("Starting server on ", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
