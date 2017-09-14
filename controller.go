package fit

import (
	"html/template"
	"net/http"
	"encoding/json"
	"fmt"
)

type ControllerInterface interface {
	Get(w *Response, r *Request, p Params)
	Post(w *Response, r *Request, p Params)
	Head(w *Response, r *Request, p Params)
	Options(w *Response, r *Request, p Params)
	Put(w *Response, r *Request, p Params)
	Patch(w *Response, r *Request, p Params)
	Delete(w *Response, r *Request, p Params)
}

// rendering Json data (© jp)
type ConterllerRenderingJson interface {
	renderingJsonAutomatically(result int, errMsg string)
	renderingJson(result int, errMsg string, datas []interface{})
}

//json result struct
type Result struct {
	Result     int  `json:"result"`// 错误码
	ErrorMsg   string `json:"errmsg"`// 错误描述
	Datas interface{} `json:"datas"`//数据
}

type Data map[interface{}]interface{}
type Controller struct {
	Data Data//界面数据
	JsonData Result
	// template data
	TplName   string
	ViewPath  string
	TplPrefix string
	TplExt    string
}

// Get adds a request function to handle GET request.
func (c *Controller) Get(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// Post adds a request function to handle post request.
func (c *Controller) Post(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// Delete adds a request function to handle delete request.
func (c *Controller) Delete(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// Put adds a request function to handle put request.
func (c *Controller) Put(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// Head adds a request function to handle head request.
func (c *Controller) Head(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// Patch adds a request function to handle patch request.
func (c *Controller) Patch(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// Options adds a request function to handle options request.
func (c *Controller) Options(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

func (c *Controller) LoadView(w *Response, tplname string) {
	var err error
	var t *template.Template
	t, err = template.ParseFiles("view"+"/" + tplname)  //从文件创建一个模板
	CheckError(err)
	err = t.Execute(w.Writer(), c.Data)
	CheckError(err)
}
// json 输出函数
func (c *Controller)ResponseToJson(w *Response)  {
	b, err := json.Marshal(c.JsonData)
	if err != nil {
		fmt.Fprint(w.Writer(), err.Error())
		return
	}
	fmt.Fprint(w.Writer(), string(b))
}
