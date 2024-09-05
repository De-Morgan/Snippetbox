package main

import (
	"html/template"
	"log"

	"github.com/morgan/snippetbox/pkg/models/database"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *database.SnippetModel
	templateCache map[string]*template.Template
}
