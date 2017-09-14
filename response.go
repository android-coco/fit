package fit

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Response struct {
	Status  int
	writer  http.ResponseWriter
	request *http.Request
	Cookie  Cookie
	data    []byte
	// defer file sending
	file      string
	redirect  string
	skipFlush bool
}

func newResponse(request *http.Request, w http.ResponseWriter) *Response {
	response := &Response{
		Status:  http.StatusOK,
		request: request,
		writer:  w,
		Cookie:  Cookie{},
	}

	return response
}

///////////////////////////////////////////////////////////////////
// Creating Responses
///////////////////////////////////////////////////////////////////

// Will produce JSON string representation of passed object,
// and send it to client
func (r *Response) Json(obj interface{}) error {
	res, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	r.writer.Header().Set("Content-Type", "application/json")
	return r.Raw(res)
}

// Will produce XML string representation of passed object,
// and send it to client
func (r *Response) Xml(obj interface{}) error {
	res, err := xml.Marshal(obj)
	if err != nil {
		return err
	}

	r.writer.Header().Set("Content-Type", "application/xml")
	return r.Raw(res)
}

// Will look for template, render it, and send rendered HTML to client.
// Second argument is data which will be passed to client.
func (r *Response) LoadView(tplName string, data interface{}) error {

	t, err := template.ParseFiles(tplName)
	if err != nil {
		r.Status = http.StatusNotFound
		return err
	}

	err = t.Execute(r.writer, data)
	if err != nil {
		r.Status = http.StatusNotFound
		return err
	}

	r.SkipFlush()
	return nil
}

// Send Raw data to client.
func (r *Response) Raw(data []byte) error {
	r.data = data
	return nil
}

// Redirect to url with status
func (r *Response) Redirect(url string) error {
	r.redirect = url
	return nil
}

// Write raw response. Implements ResponseWriter.Write.
func (r *Response) Write(b []byte) (int, error) {
	return r.writer.Write(b)
}

// Get Header. Implements ResponseWriter.Header.
func (r *Response) Header() http.Header {
	return r.writer.Header()
}

// Write Header. Implements ResponseWriter.WriterHeader.
func (r *Response) WriteHeader(s int) {
	r.writer.WriteHeader(s)
}

// Get http.ResponseWriter directly
func (r *Response) Writer() http.ResponseWriter {
	return r.writer
}

// Checking if file exist.
// todo: consider moving this to utils.go
func (r *Response) fileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}

	r.Status = http.StatusNotFound

	return false
}

// Find file, and send it to client.
func (r *Response) File(path string) error {
	abspath, err := filepath.Abs(path)

	if err != nil {
		return err
	}

	if !r.fileExists(abspath) {
		return err
	}

	r.file = abspath
	return nil
}

// Serving static file.
func (r *Response) serveFile(file string) error {

	http.ServeFile(r.writer, r.request, file)
	return nil
}

// Will be called from ``flush`` Response method if user called ``File`` method.
func (r *Response) sendFile() {

	base := filepath.Base(r.file)
	r.writer.Header().Set("Content-Disposition", "attachment; filename="+base)
	http.ServeFile(r.writer, r.request, r.file)
}

///////////////////////////////////////////////////////////////////
// Writing Response
///////////////////////////////////////////////////////////////////

func (r *Response) SkipFlush() {
	r.skipFlush = true
}

// Write result to ResponseWriter.
func (r *Response) flush() {
	if r.skipFlush {

		return
	}

	// set all cookies to response object
	for _, v := range r.Cookie {
		//fmt.Printf("k == %v\n", k)
		//fmt.Printf("v == %v\n", v)

		http.SetCookie(r.writer, v)
	}

	//r.Cookie.SetResponseCookie(r.writer)

	// in case of file call separate function for piping file to client
	if len(r.file) > 0 {
		r.sendFile()
	} else if len(r.redirect) > 0 {
		http.Redirect(r.writer, r.request, r.redirect, r.Status)
	} else if r.data != nil && len(r.data) > 0 {
		r.writer.WriteHeader(r.Status)
		r.writer.Write(r.data)
	}
}
