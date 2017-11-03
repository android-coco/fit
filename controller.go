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
	RenderingJsonAutomatically(result int, errMsg string)
	RenderingJson(result int, errMsg string, datas interface{})
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
// 重定向
func (c *Controller)Redirect(w *Response, r *Request,url string, code int){
	http.Redirect(w.writer,r.Request,url,code)
}

func (c *Controller) LoadView(w *Response, tplname ...string) {
	var err error
	var t *template.Template
	var args []string
	for _, arg := range tplname {
		arg = "view/" + arg
		args = append(args, arg)
	}

	t, err = template.ParseFiles(args...)  //从文件创建一个模板

    if err != nil {
        Logger().LogError("Fatal error ", err.Error())
    }

	err = t.Execute(w.Writer(), c.Data)
	
    if err != nil {
        Logger().LogError("Fatal error ", err.Error())
    }
}

/*Checking session and then load a special view */
func (c *Controller) LoadViewSafely(w *Response, r *Request, tplname ...string) (success bool) {
	session, err_s := GlobalManager().SessionStart(w, r)
	if err_s != nil || session == nil {
		// Session失效 重新登录
		Logger().LogError("JP ", "Session失效 重新登录" + err_s.Error())
		return false
	} else {
		userinfo := session.Get("UserInfo")
		if userinfo == nil {
			// 未登录
			Logger().LogDebug("JP ", "未登录")
			return false
		} else {
			// 已登录
			c.LoadView(w,tplname...)
			Logger().LogDebug("JP ", "已登录")
			return true
		}
	}
}



// json 输出函数
func (c *Controller)ResponseToJson(w *Response)  {
	if c.JsonData.Datas == nil{
		c.JsonData.Datas = []interface{}{}
	}
	b, err := json.Marshal(c.JsonData)
	if err != nil {
		fmt.Fprint(w.Writer(), err.Error())
		return
	}
	fmt.Fprint(w.Writer(), string(b))
}
