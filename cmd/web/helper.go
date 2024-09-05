package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// The serverError writes an error message and stack trace to the errorLogger
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding descri
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, data *templateData) {

	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("template doees not exist: %s", name))
		return
	}
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the  http.ResponseWriter.
	//to catch any error present
	err := ts.Execute(buf, app.addDefaultData(data, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, _ *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td
}
